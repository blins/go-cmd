/*
Пакет расширяет минимальный функционал стандартного пакета flag
*/
package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

// интерфейс, который должна реализовывать команда
type Command interface {
	// возвращает FlagSet для данной команды
	GetFlags() *flag.FlagSet
	// Разбор параметров. Результатом должно стать остаток массива args, которые не относится к данной команде
	ParseArgs(args []string) ([]string, error)
	// запустить выполнение команды. Если команда асинхронна, то возвращается Waiter (можно использовать sync.WaitGroup), иначе Waiter = nil
	Run(ctx context.Context) (Waiter, error)
}

// фабрика команд
type CommandFabric interface {
	Create() Command
}

// фабрика команд ввиде функции
type CommandFabricFunc func() Command

func (f CommandFabricFunc) Create() Command {
	return f()
}

// зарегистрировать фабрику команд
func RegisterFabric(name string, fabric CommandFabric) {
	if cmdRegistry == nil {
		cmdRegistry = make(map[string]CommandFabric)
	}
	cmdRegistry[name] = fabric
}

// взять фабрику команд
func GetFabric(name string) (CommandFabric, bool) {
	if cmdRegistry == nil {
		return nil, false
	}
	v, ok := cmdRegistry[name]
	return v, ok
}

var (
	cmdRegistry map[string]CommandFabric
	wg          *sync.WaitGroup
)

type Waiter interface {
	Wait()
}

// запустить парсинг параметров команд и последовательный запуск их
func ParseAndRun(defaultCmds []string, ctx context.Context) Waiter {
	wg = &sync.WaitGroup{}
	args := flag.Args()
	if len(args) == 0 {
		args = defaultCmds
	}

	for len(args) > 0 {
		cmdName := args[0]
		if cmdf, ok := GetFabric(cmdName); ok {
			cmd := cmdf.Create()
			args, _ = cmd.ParseArgs(args[1:])
			wg.Add(1)
			go func(c Command) {
				wait, err := c.Run(ctx)
				if err != nil {
					log.Printf("Error on %v: %v\n", c.GetFlags().Name(), err)
				}
				if wait != nil {
					wait.Wait()
				}
				wg.Done()
			}(cmd)
		} else {
			log.Fatalln("unknown command")
		}
	}
	return wg
}

func Usage() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "%s [global_options] [command [command_options]]...:\n", os.Args[0])
	flag.PrintDefaults()
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Commands:\n")
	for _, cmdf := range cmdRegistry {
		cmd := cmdf.Create()
		cmd.GetFlags().Usage()
	}
}

func init() {
	flag.Usage = Usage
}
