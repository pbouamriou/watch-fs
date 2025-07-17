package ui

import (
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jesseduffield/gocui"
)

// Events gère la logique des événements du watcher et leur traitement
// (ajout, filtrage, tri, agrégation)
type Events struct {
	ui *UI
}

func NewEvents(ui *UI) *Events {
	return &Events{ui: ui}
}

// addEvent ajoute un nouvel événement à l’état
func (e *Events) addEvent(path string, operation fsnotify.Op, isDir bool) {
	if e.ui.state.AggregateEvents {
		// Vérifie si un événement similaire existe dans la dernière seconde
		now := time.Now()
		for _, event := range e.ui.state.Events {
			if event.Path == path && event.Operation == operation &&
				now.Sub(event.Timestamp) < time.Second {
				event.Count++
				event.Timestamp = now
				return
			}
		}
	}

	// Ajoute un nouvel événement
	event := &FileEvent{
		Path:      path,
		Operation: operation,
		Timestamp: time.Now(),
		IsDir:     isDir,
		Count:     1,
	}
	e.ui.state.Events = append(e.ui.state.Events, event)

	// Limite le nombre d’événements
	if len(e.ui.state.Events) > e.ui.state.MaxEvents {
		e.ui.state.Events = e.ui.state.Events[1:]
	}

	// Met à jour l’UI si initialisée
	if e.ui.gui != nil {
		e.ui.gui.Update(func(g *gocui.Gui) error {
			if v, err := g.View(EventsView); err == nil {
				e.ui.views.UpdateEventsView(v)
			}
			if v, err := g.View(StatusView); err == nil {
				e.ui.views.UpdateStatusView(v)
			}
			if v, err := g.View(FilterView); err == nil {
				e.ui.views.UpdateFilterView(v)
			}
			return nil
		})
	}
}

// watchEvents écoute les événements du watcher
func (e *Events) watchEvents() {
	for {
		select {
		case event, ok := <-e.ui.watcher.Events():
			if !ok {
				return
			}
			info, err := e.ui.getFileInfo(event.Name)
			isDir := false
			if err == nil {
				isDir = info.IsDir()
			}
			e.addEvent(event.Name, event.Op, isDir)
		case err, ok := <-e.ui.watcher.Errors():
			if !ok {
				return
			}
			// Ignore l’erreur pour l’instant
			_ = err
		}
	}
}

// getFilteredEvents retourne les événements filtrés et triés
func (e *Events) getFilteredEvents() []*FileEvent {
	filtered := make([]*FileEvent, 0)
	for _, event := range e.ui.state.Events {
		// Filtre chemin
		if e.ui.state.Filter.PathFilter != "" &&
			!strings.Contains(strings.ToLower(event.Path), strings.ToLower(e.ui.state.Filter.PathFilter)) {
			continue
		}
		// Filtre opération
		if e.ui.state.Filter.OperationFilter != 0 && event.Operation != e.ui.state.Filter.OperationFilter {
			continue
		}
		// Filtre type
		if event.IsDir && !e.ui.state.Filter.ShowDirs {
			continue
		}
		if !event.IsDir && !e.ui.state.Filter.ShowFiles {
			continue
		}
		filtered = append(filtered, event)
	}
	e.sortEvents(filtered)
	return filtered
}

// sortEvents trie les événements selon l’option courante
func (e *Events) sortEvents(events []*FileEvent) {
	switch e.ui.state.SortOption {
	case SortByTime:
		sort.Slice(events, func(i, j int) bool {
			return events[i].Timestamp.After(events[j].Timestamp)
		})
	case SortByPath:
		sort.Slice(events, func(i, j int) bool {
			return events[i].Path < events[j].Path
		})
	case SortByOperation:
		sort.Slice(events, func(i, j int) bool {
			return events[i].Operation < events[j].Operation
		})
	case SortByCount:
		sort.Slice(events, func(i, j int) bool {
			return events[i].Count > events[j].Count
		})
	}
}

// getSortOptionName retourne le nom de l’option de tri courante
func (e *Events) getSortOptionName() string {
	switch e.ui.state.SortOption {
	case SortByTime:
		return "Time"
	case SortByPath:
		return "Path"
	case SortByOperation:
		return "Operation"
	case SortByCount:
		return "Count"
	default:
		return "Unknown"
	}
}
