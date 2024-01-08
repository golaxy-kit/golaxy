package golaxy

import (
	"kit.golaxy.org/golaxy/runtime"
	"time"
)

func (rt *RuntimeBehavior) loopingWithFrame() {
	gcTicker := time.NewTicker(rt.opts.GCInterval)
	defer gcTicker.Stop()

	frame := runtime.UnsafeFrame(rt.opts.Frame)
	go rt.makeFrameTasks(frame.GetCurFrames()+1, frame.GetTotalFrames(), frame.GetTargetFPS())

loop:
	for rt.frameLoopBegin(); ; {
		select {
		case task, ok := <-rt.processQueue:
			if !ok {
				break loop
			}
			rt.runTask(task)

		case <-gcTicker.C:
			rt.runGC()

		case <-rt.ctx.Done():
			break loop
		}
	}

	close(rt.processQueue)

loopEnding:
	for {
		select {
		case task, ok := <-rt.processQueue:
			if !ok {
				break loopEnding
			}
			rt.runTask(task)

		default:
			break loopEnding
		}
	}

	rt.frameLoopEnd()
}

func (rt *RuntimeBehavior) makeFrameTasks(curFrames, totalFrames uint64, targetFPS float32) {
	updateTicker := time.NewTicker(time.Duration(float64(time.Second) / float64(targetFPS)))
	defer updateTicker.Stop()

	for {
		if totalFrames > 0 && curFrames >= totalFrames {
			rt.Terminate()
			return
		}

		select {
		case <-updateTicker.C:
			func() {
				defer func() {
					recover()
				}()
				select {
				case rt.processQueue <- _Task{kind: _TaskKind_Frame, action: rt.frameLoop}:
					curFrames++
				case <-rt.ctx.Done():
				}
			}()
		case <-rt.ctx.Done():
			return
		}
	}
}

func (rt *RuntimeBehavior) frameLoop(...any) {
	rt.frameLoopEnd()
	rt.frameLoopBegin()
}

func (rt *RuntimeBehavior) frameLoopBegin() {
	rt.changeRunningState(runtime.RunningState_FrameLoopBegin)
	rt.changeRunningState(runtime.RunningState_FrameUpdateBegin)

	emitEventUpdate(&rt.eventUpdate)
	emitEventLateUpdate(&rt.eventLateUpdate)

	rt.changeRunningState(runtime.RunningState_FrameUpdateEnd)
}

func (rt *RuntimeBehavior) frameLoopEnd() {
	rt.changeRunningState(runtime.RunningState_FrameLoopEnd)

	frame := runtime.UnsafeFrame(rt.opts.Frame)
	frame.SetCurFrames(frame.GetCurFrames() + 1)
}
