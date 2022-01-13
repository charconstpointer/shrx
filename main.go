package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	//foo2
	log1 := &Log{
		Messages: []string{"a", "b", "c"},
	}
	log2 := &Log{
		Messages: []string{"a", "b", "c"},
	}
	part1 := &MemPart{
		index:  0,
		Stream: log1,
	}
	part2 := &MemPart{
		index:  1,
		Stream: log2,
	}
	var g errgroup.Group
	g.Go(func() error {
		return part1.Run()
	})
	g.Go(func() error {
		return part2.Run()
	})
	keeper := NewKeeper(part1, part2)
	keeper2 := NewKeeper(keeper)
	go func(k *Keeper) {
		for {
			k.Write("y")
			time.Sleep(time.Millisecond * 10)
		}
	}(keeper2)
	time.Sleep(time.Second * 10)
}

type Partition interface {
	Write(msg string)
}

type MemPart struct {
	Stream *Log
	index  int
	mu     sync.Mutex
}

func (p *MemPart) Run() error {
	for msg := p.Stream.Next(); msg != nil; msg = p.Stream.Next() {
		fmt.Printf("%d: %s\n", p.index, *msg)
		time.Sleep(time.Millisecond * 99)
	}
	return nil
}

func (p *MemPart) Write(msg string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Stream.Append(msg)
}

type Log struct {
	Messages []string
}

func (l *Log) Next() *string {
	if len(l.Messages) == 0 {
		return nil
	}
	msg := l.Messages[0]
	l.Messages = l.Messages[1:]
	return &msg
}

func (l *Log) Append(msg string) {
	l.Messages = append(l.Messages, msg)
}

type Keeper struct {
	Partitions []Partition
	next       int
}

func NewKeeper(parts ...Partition) *Keeper {
	if len(parts) == 0 {
		return nil
	}
	return &Keeper{
		Partitions: parts,
		next:       0,
	}
}

func (k *Keeper) Write(msg string) {
	if len(k.Partitions) == 0 {
		return
	}
	k.Partitions[k.next].Write(msg)
	k.next = (k.next + 1) % len(k.Partitions)
}
