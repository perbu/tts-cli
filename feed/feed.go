package feed

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

type abtractIntelligence interface {
	Summary(ctx context.Context, content string) (string, error)
	Illustration(ctx context.Context, content string) ([]byte, error)
	Speech(ctx context.Context, content string) (io.ReadCloser, error)
}

type FeedManager struct {
	openAI   abtractIntelligence
	episodes []Episode
	logger   *slog.Logger
}

func New(openAI abtractIntelligence, logger *slog.Logger) *FeedManager {
	return &FeedManager{
		episodes: make([]Episode, 0),
		openAI:   openAI,
		logger:   logger,
	}
}
func (f *FeedManager) Scan(ctx context.Context, directory string) error {
	logger := f.logger.With("function", "scan", "directory", directory)
	// Scan the directory for txt files
	fh, err := os.Open(directory)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer fh.Close()

	files, err := fh.Readdir(-1)
	if err != nil {
		return fmt.Errorf("Readdir: %w", err)
	}
	logger.Debug("files", "files", len(files))

	// Sort files by creation time
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	// For each file, add an episode
	for _, file := range files {
		logger.Debug("file", "name", file.Name(), "isDir", file.IsDir())
		if file.IsDir() {
			continue
		}
		if path.Ext(file.Name()) != ".txt" {
			logger.Debug("skipping non-txt file")
			continue
		}
		// skip the file if the name contains more than two dots:
		if strings.HasSuffix(file.Name(), ".summary.txt") {
			logger.Debug("skipping summary")
			continue
		}
		if err := f.AddEpisode(ctx, file.Name()); err != nil {
			return fmt.Errorf("AddEpisode: %w", err)
		}
	}

	return nil
}

func (f *FeedManager) AddEpisode(ctx context.Context, contentFile string) error {
	logger := f.logger.With("contentFile", contentFile)
	logger.Debug("AddEpisode", "contentFile", contentFile)
	if contentFile == "" {
		return fmt.Errorf("empty contentFile")
	}
	// Parse the contentFile to get the content, Summary, and illustration filenames
	// Add the episode to the feed
	e := Episode{
		ContentFile:      contentFile,
		summaryFile:      contentFile + ".summary.txt",
		IllustrationFile: contentFile + ".png",
		AudioFile:        contentFile + ".mp3",
	}

	// get the created and updated times of the content file.
	stat, err := os.Stat(e.ContentFile)
	if err != nil {
		return fmt.Errorf("os.Stat(%s): %w", e.ContentFile, err)
	}
	logger.Debug("stat", "stat", stat)
	e.Created = stat.ModTime()

	contentBytes, err := os.ReadFile(e.ContentFile)
	if err != nil {
		return fmt.Errorf("os.ReadFile(%s): %w", e.ContentFile, err)
	}
	logger.Debug("content", "len", len(contentBytes))
	if len(contentBytes) == 0 {
		return fmt.Errorf("empty content")
	}
	e.content = string(contentBytes)

	summaryBytes, err := os.ReadFile(e.summaryFile)
	e.Summary = string(summaryBytes)
	if err != nil {
		logger.Info("no summary file found, generating one...")
		e.Summary, err = f.openAI.Summary(ctx, e.content)
		if err != nil {
			return fmt.Errorf("ai.Summary: %w", err)
		}
		if err := os.WriteFile(e.summaryFile, []byte(e.Summary), 0644); err != nil {
			return fmt.Errorf("os.WriteFile: %w", err)
		}
		logger.Debug("summary written", "len", len(e.Summary))
	}

	if _, err := os.Stat(e.IllustrationFile); err != nil {
		logger.Info("no illustration file found, generating one...")
		// try to generate illustration
		illustrationBytes, err := f.openAI.Illustration(ctx, e.Summary)
		if err != nil {
			return fmt.Errorf("ai.Illustration: %w", err)
		}
		logger.Debug("got illustration bytes back", "len", len(illustrationBytes))
		if err := os.WriteFile(e.IllustrationFile, illustrationBytes, 0644); err != nil {
			return fmt.Errorf("os.WriteFile: %w", err)
		}
		logger.Debug("illustration written", "len", len(illustrationBytes))
	}

	// check if the audio file exists
	stat, err = os.Stat(e.AudioFile)
	if err != nil {
		logger.Info("no audio file found, generating one...")
		// try to generate audio
		wr, err := f.openAI.Speech(ctx, e.content)
		if err != nil {
			return fmt.Errorf("ai.Speech: %w", err)
		}
		fh, err := os.OpenFile(e.AudioFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("os.OpenFile(%s): %w", e.AudioFile, err)
		}
		n, err := io.Copy(fh, wr)
		if err != nil {
			return fmt.Errorf("io.Copy: %w", err)
		}
		logger.Info("audio written", "len", n)
		err = fh.Close()
		if err != nil {
			return fmt.Errorf("fh.Close: %w", err)
		}
		// re-stat the file:
		stat, err = os.Stat(e.AudioFile)
		if err != nil {
			return fmt.Errorf("os.Stat(%s): %w", e.AudioFile, err)
		}
	}
	// get the updated time of the audio file and use it as the updated time of the episode.
	e.Updated = stat.ModTime()

	// things seems ok here, add the episode to the feed.
	logger.Debug("adding episode to feed", "feed length", len(f.episodes))
	f.episodes = append(f.episodes, e)
	return nil
}

func (f *FeedManager) GenerateRSS() (string, error) {
	channel := &Channel{
		Title:         "Your Podcast Title",
		Link:          "https://yourpodcastwebsite.com",
		Description:   "Your podcast description",
		Language:      "en-us",
		PubDate:       time.Now().Format(time.RFC1123Z),
		LastBuildDate: time.Now().Format(time.RFC1123Z),
	}

	for _, episode := range f.episodes {
		item := &Item{
			Title:       episode.ContentFile,
			Link:        "https://yourpodcastwebsite.com/" + episode.ContentFile,
			Description: episode.Summary,
			PubDate:     episode.Created.Format(time.RFC1123Z),
			Guid:        "https://yourpodcastwebsite.com/" + episode.ContentFile,
		}

		if episode.AudioFile != "" {
			item.Enclosure = &Enclosure{
				URL:    "https://yourpodcastwebsite.com/" + episode.AudioFile,
				Type:   "audio/mpeg",
				Length: "0", // You should set the actual file size here
			}
		}

		if episode.IllustrationFile != "" {
			item.ITunesImage = &ITunesImage{
				Href: "https://yourpodcastwebsite.com/" + episode.IllustrationFile,
			}
		}

		channel.Items = append(channel.Items, item)
	}

	rss := &RSS{
		Version: "2.0",
		Channel: channel,
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(output), nil
}
