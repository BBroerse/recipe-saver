package events

import (
	"context"
	"fmt"
	"recipe-processor/internal/shared/logger"
	"sync"
	"time"
)

const (
	// DefaultWorkerCount is the default number of worker goroutines
	DefaultWorkerCount = 10
	// DefaultChannelBuffer is the default buffer size for event channels
	DefaultChannelBuffer = 100
	// DefaultHandlerTimeout is the default timeout for event handlers
	DefaultHandlerTimeout = 30 * time.Second
)

// MemoryEventBus is an in-memory event bus using channels and goroutines
type MemoryEventBus struct {
	handlers       map[string][]EventHandler
	eventQueue     chan Event
	workerCount    int
	handlerTimeout time.Duration
	logger         logger.Logger
	mu             sync.RWMutex
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
}

// Config holds configuration for the memory event bus
type Config struct {
	WorkerCount    int
	ChannelBuffer  int
	HandlerTimeout time.Duration
}

// NewMemoryEventBus creates a new in-memory event bus with default config
func NewMemoryEventBus(log logger.Logger) EventBus {
	return NewMemoryEventBusWithConfig(log, Config{
		WorkerCount:    DefaultWorkerCount,
		ChannelBuffer:  DefaultChannelBuffer,
		HandlerTimeout: DefaultHandlerTimeout,
	})
}

// NewMemoryEventBusWithConfig creates a new in-memory event bus with custom config
func NewMemoryEventBusWithConfig(log logger.Logger, cfg Config) EventBus {
	return &MemoryEventBus{
		handlers:       make(map[string][]EventHandler),
		eventQueue:     make(chan Event, cfg.ChannelBuffer),
		workerCount:    cfg.WorkerCount,
		handlerTimeout: cfg.HandlerTimeout,
		logger:         log,
	}
}

// Start begins processing events with worker goroutines
func (eb *MemoryEventBus) Start(ctx context.Context) error {
	eb.ctx, eb.cancel = context.WithCancel(ctx)

	// Start worker pool
	for i := 0; i < eb.workerCount; i++ {
		eb.wg.Add(1)
		go eb.worker(i)
	}

	eb.logger.Info("Event bus started",
		logger.Int("workers", eb.workerCount),
		logger.Int("buffer_size", cap(eb.eventQueue)),
	)
	return nil
}

// Stop gracefully shuts down the event bus
func (eb *MemoryEventBus) Stop() error {
	if eb.cancel != nil {
		eb.cancel()
	}

	// Close event queue to signal workers to stop
	close(eb.eventQueue)

	// Wait for all workers to finish processing
	eb.wg.Wait()

	eb.logger.Info("Event bus stopped")
	return nil
}

// Publish sends an event to all registered handlers
func (eb *MemoryEventBus) Publish(ctx context.Context, event Event) error {
	select {
	case eb.eventQueue <- event:
		eb.logger.Debug("Event published",
			logger.String("event_type", event.EventType()),
		)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("publish cancelled: %w", ctx.Err())
	case <-eb.ctx.Done():
		return fmt.Errorf("event bus stopped: %w", eb.ctx.Err())
	}
}

// Subscribe registers a handler for a specific event type
func (eb *MemoryEventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	eb.logger.Info("Handler subscribed",
		logger.String("event_type", eventType),
		logger.Int("total_handlers", len(eb.handlers[eventType])),
	)
}

// worker processes events from the queue
func (eb *MemoryEventBus) worker(id int) {
	defer eb.wg.Done()

	eb.logger.Debug("Worker started", logger.Int("worker_id", id))

	for {
		select {
		case event, ok := <-eb.eventQueue:
			if !ok {
				// Channel closed, exit worker
				eb.logger.Debug("Worker stopped", logger.Int("worker_id", id))
				return
			}

			// Process event
			eb.processEvent(event)

		case <-eb.ctx.Done():
			eb.logger.Debug("Worker cancelled", logger.Int("worker_id", id))
			return
		}
	}
}

// processEvent dispatches an event to all registered handlers
func (eb *MemoryEventBus) processEvent(event Event) {
	eb.mu.RLock()
	handlers := eb.handlers[event.EventType()]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		eb.logger.Warn("No handlers registered for event",
			logger.String("event_type", event.EventType()),
		)
		return
	}

	// Execute all handlers for this event type
	for i, handler := range handlers {
		// Create context with timeout for handler execution
		handlerCtx, cancel := context.WithTimeout(eb.ctx, eb.handlerTimeout)

		start := time.Now()
		err := handler(handlerCtx, event)
		duration := time.Since(start)

		if err != nil {
			eb.logger.Error("Handler failed",
				logger.String("event_type", event.EventType()),
				logger.Int("handler_index", i),
				logger.Duration("duration", duration),
				logger.Error(err),
			)
		} else {
			eb.logger.Debug("Handler succeeded",
				logger.String("event_type", event.EventType()),
				logger.Int("handler_index", i),
				logger.Duration("duration", duration),
			)
		}

		cancel()
	}
}
