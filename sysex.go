package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
#include <stdio.h>
#include <unistd.h>

extern void sysex_callback(MIDISysexSendRequest *request);

static inline void MIDISysexProc(MIDISysexSendRequest *request)
{
    sysex_callback(request);
}

typedef void (*midi_sysex_proc)(MIDISysexSendRequest *request);

static midi_sysex_proc getSysexProc()
{
    return *MIDISysexProc;
}

*/
import "C"
import (
	"errors"
	"unsafe"
)

import "fmt"

//export sysex_callback
func sysex_callback(r *C.MIDISysexSendRequest) {
	foo := *(*SysexMessage)(r.completionRefCon)
	foo.SysexProc(&foo)
}

type SysexProc func(request *SysexMessage)

type SysexMessage struct {
	SysexsendRequest C.MIDISysexSendRequest
	SysexProc        SysexProc
	Message          []byte
}

func NewSysexMessage(destination *Destination, data []byte, sysexProc SysexProc) SysexMessage {

	var sysexMessage = SysexMessage{SysexProc: sysexProc, Message: data}
	var sendRequest C.MIDISysexSendRequest

	sendRequest.destination = destination.endpoint
	sendRequest.data = (*C.Byte)(unsafe.Pointer(&data[0]))
	sendRequest.bytesToSend = (C.UInt32)(len(data))
	sendRequest.completionProc = (C.MIDICompletionProc)(C.getSysexProc())
	sendRequest.completionRefCon = unsafe.Pointer(&sysexMessage)

	sysexMessage.SysexsendRequest = sendRequest

	return sysexMessage

}

func (sysex *SysexMessage) Send() (err error) {

	osStatus := C.MIDISendSysex(&sysex.SysexsendRequest)

	if osStatus != C.noErr {
		err = errors.New(fmt.Sprintf("%d: failed to send Sysex", int(osStatus)))
	}

	return
}
