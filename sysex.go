package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
#include <stdio.h>
#include <unistd.h>

static void MIDISysexProc(MIDISysexSendRequest *request)
{


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

/*
// structs
struct MIDISysexSendRequest {
    MIDIEndpointRef destination;        // destination.endpoint
    const Byte *data;                   // data
    UInt32 bytesToSend;                 // C.ByteCount(len(p)
    Boolean complete;
    Byte reserved[3];
    MIDICompletionProc completionProc;  // completionProc
    void *completionRefCon;             //
};

// functions
extern OSStatus MIDISendSysex(
    MIDISysexSendRequest *request );

*/

type SysexProc func(request *SysexMessage)

type SysexMessage struct {
	sysexSendRequest C.MIDISysexSendRequest
	SysexProc        SysexProc
}

func MyCallback(x C.MIDISysexSendRequest) {
	fmt.Println("callback with", x)
}

func NewSysexMessage(destination *Destination, data []byte, sysexProc SysexProc) (sysexMessage SysexMessage) {

	var SysexRequest C.MIDISysexSendRequest
	SysexRequest.destination = destination.endpoint
	SysexRequest.data = (*C.Byte)(unsafe.Pointer(&data[0]))
	SysexRequest.bytesToSend = (C.UInt32)(len(data))
	SysexRequest.completionProc = (C.MIDICompletionProc)(C.getSysexProc())

	sysexMessage = SysexMessage{SysexRequest, sysexProc}

	return

}

/*


func NewInputPort(client Client, name string, readProc ReadProc) (inputPort InputPort, err error) {
    var port C.MIDIPortRef

    stringToCFString(name, func(cfName C.CFStringRef) {
        osStatus := C.MIDIInputPortCreate(client.client,
            cfName,
            (C.MIDIReadProc)(C.getProc()),
            unsafe.Pointer(uintptr(0)),
            &port)

        if osStatus != C.noErr {
            err = errors.New(fmt.Sprintf("%d: failed to create a port", int(osStatus)))
        } else {
            inputPort = InputPort{port, readProc, make([]*C.int, 0)}
        }
    })

    return
}

*/

func (sysex *SysexMessage) Send() (err error) {

	osStatus := C.MIDISendSysex(&sysex.sysexSendRequest)

	if osStatus != C.noErr {
		err = errors.New(fmt.Sprintf("%d: failed to send Sysex", int(osStatus)))
	}

	return
}
