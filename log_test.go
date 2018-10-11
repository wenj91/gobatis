package gobatis

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestInfo(t *testing.T) {
	SetLevel(LOG_LEVEL_DEBUG)
	Info("test info -> level debug")

	SetLevel(LOG_LEVEL_OFF)
	Info("test info -> level off")
}

func TestFileLogger_Info(t *testing.T)  {
	logger := NewFileLog("d:/logs/nohup.out", LOG_LEVEL_DEBUG)
	logger.Info("test file logger info")
}

func TestFileLogger(t *testing.T) {
	fileName := "d:/logs/nohup.out"
	fileName = fileName

	// open a file
	logger := NewFileLog(fileName, LOG_LEVEL_DEBUG)

	now := time.Now()
	var w sync.WaitGroup
	w.Add(1000 * 3)
	for i := 0; i < 1000; i++ {
		go func() {
			logger.Info("test happy time--info")
			w.Done()
		}()
		go func() {
			logger.Debug("test happy time--debug")
			w.Done()
		}()

		go func() {
			logger.Fatal("test happy time--fatal")
			w.Done()
		}()
	}

	w.Wait()
	cur := time.Now().Sub(now).Nanoseconds()
	fmt.Println("cur:", cur)

	w.Add(1000 * 3)
	now = time.Now()
	for i := 0; i < 1000; i++ {
		go func() {
			Info("test happy time--info")
			w.Done()
		}()
		go func() {
			Debug("test happy time--debug")
			w.Done()
		}()

		go func() {
			Fatal("test happy time--fatal")
			w.Done()
		}()
	}

	w.Wait()
	cur = time.Now().Sub(now).Nanoseconds()
	fmt.Println("cur:", cur)

	w.Add(1000*3)
	now = time.Now()
	for i := 0; i < 1000; i++ {
		go func() {
			log.Println("test happy time--info")
			w.Done()
		}()
		go func() {
			log.Println("test happy time--debug")
			w.Done()
		}()

		go func() {
			log.Println("test happy time--fatal")
			w.Done()
		}()
	}

	w.Wait()
	cur = time.Now().Sub(now).Nanoseconds()
	fmt.Println("cur:", cur)

	fmt.Scanf("%d", nil)
}
