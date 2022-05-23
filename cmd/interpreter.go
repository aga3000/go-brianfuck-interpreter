package main

import (
	"bufio"
	"fmt"
	"github.com/aga3000/go-brianfuck"
	"github.com/pkg/errors"
	"io"
	"log"
	"math"
	"os"
)

type Uint8WriterWrapper struct {
	Writer io.Writer
}

func (w Uint8WriterWrapper) WriteByte(c byte) error {
	_, err := fmt.Fprintf(w.Writer, "[%d] ", c)
	return err
}

func main() {
	var sourceReader io.RuneReader
	reader := bufio.NewReader(os.Stdin)
	writer := Uint8WriterWrapper{Writer: os.Stdout}
	if len(os.Args) == 1 {
		sourceReader = bufio.NewReader(os.Stdin)
	} else if len(os.Args) == 2  {
		sourcePath := os.Args[1]
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			log.Fatalf("failed to read source file: %v", err)
		}
		sourceReader = bufio.NewReader(sourceFile)
	} else {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s [source]", os.Args[0])
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
	cmdMap := brainfuck.GetDefaultCommandMap()
	cmdMap['?'] = func(state brainfuck.InterpreterState) error {
		cellValue := state.Memory.Read()
		sqrtCellValueF := math.Sqrt(float64(cellValue))
		sqrtCellValue := brainfuck.Cell(math.Ceil(sqrtCellValueF))
		state.Memory.Write(sqrtCellValue)
		return nil
	}
	runner, err := brainfuck.NewInterpreterRunner(
		reader, writer,
		brainfuck.WithCommands(cmdMap),
		brainfuck.WithUnknownCharPolicy(brainfuck.IgnoreUnknownCharsPolicy),
	)
	if err != nil {
		log.Fatalf("failed to execute source: %+v", err)
	}
	for charPos := 0; true; charPos++ {
		ch, _, codeReaderErr := sourceReader.ReadRune()
		if errors.Is(codeReaderErr, io.EOF) {
			break
		}
		err = runner.Execute(brainfuck.CommandChar(ch))
		if err != nil {
			if errors.Is(err, brainfuck.EOR) {
				break
			}
			log.Fatalf("execution failed at position %d on command %c: %+v", charPos, ch, err)
		}
	}
}
