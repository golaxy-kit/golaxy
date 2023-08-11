// Code generated by "stringer -type RunningState"; DO NOT EDIT.

package runtime

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RunningState_Birth-0]
	_ = x[RunningState_Starting-1]
	_ = x[RunningState_Started-2]
	_ = x[RunningState_FrameLoopBegin-3]
	_ = x[RunningState_FrameUpdateBegin-4]
	_ = x[RunningState_FrameUpdateEnd-5]
	_ = x[RunningState_FrameLoopEnd-6]
	_ = x[RunningState_AsyncProcessingBegin-7]
	_ = x[RunningState_AsyncProcessingEnd-8]
	_ = x[RunningState_Terminating-9]
	_ = x[RunningState_Terminated-10]
}

const _RunningState_name = "RunningState_BirthRunningState_StartingRunningState_StartedRunningState_FrameLoopBeginRunningState_FrameUpdateBeginRunningState_FrameUpdateEndRunningState_FrameLoopEndRunningState_AsyncProcessingBeginRunningState_AsyncProcessingEndRunningState_TerminatingRunningState_Terminated"

var _RunningState_index = [...]uint16{0, 18, 39, 59, 86, 115, 142, 167, 200, 231, 255, 278}

func (i RunningState) String() string {
	if i < 0 || i >= RunningState(len(_RunningState_index)-1) {
		return "RunningState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RunningState_name[_RunningState_index[i]:_RunningState_index[i+1]]
}
