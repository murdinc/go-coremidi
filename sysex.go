package coremidi

/*
#cgo LDFLAGS: -framework CoreMIDI -framework CoreFoundation
#include <CoreMIDI/CoreMIDI.h>
#include <stdio.h>
#include <unistd.h>


extern void sysex_callback(MIDISysexSendRequest *prequest);

static inline void MIDISysexProc(MIDISysexSendRequest *request)
{

    printf("(c )MIDISysexProcdata: %x\n\n", request->data);
    printf("(c) MIDISysexProc complete: %d\n\n", request->complete);
    sysex_callback(request);


}

typedef void (*midi_sysex_proc)(MIDISysexSendRequest *request);

static midi_sysex_proc getSysexProc()
{
    print("(c) midi_sysex_proc");
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
func sysex_callback(p1 *C.MIDISysexSendRequest) {
	fmt.Println("(go) sysex_callback")
	fmt.Printf("Data: %X ( should be: 0xF0, 0x43, 0x20, 0x09, 0xF7 )", p1.data)

	foo := *(*func(*C.MIDISysexSendRequest))(p1.completionRefCon)
	foo(p1)
}

/*
// structs
struct MIDISysexSendRequest {
    MIDIEndpointRef destination;        // destination.endpoint
    const Byte *data;                   // data
    UInt32 bytesToSend;                 // C.ByteCount(len(p)
    Boolean complete;
    Byte reserved[3];
    MIDICompletionProc completionProc;  // completionProc
    void *completionRefCon;             // name?
};

// functions
extern OSStatus MIDISendSysex(
    MIDISysexSendRequest *request );

typedef void ( *MIDICompletionProc)(
    MIDISysexSendRequest *request);

*/

type SysexProc func(request *SysexMessage)

type SysexMessage struct {
	sysexSendRequest C.MIDISysexSendRequest
	SysexProc        SysexProc
}

func MyCallback(x *C.MIDISysexSendRequest) {
	fmt.Println("MyCallback")
	fmt.Println("callback with", x)

}

// we store it in a global variable so that the garbage collector
// doesn't clean up the memory for any temporary variables created.
var MyCallbackFunc = MyCallback

func NewSysexMessage(destination *Destination, data []byte, sysexProc SysexProc) (sysexMessage SysexMessage) {
	var SysexRequest C.MIDISysexSendRequest

	stringToCFString("test", func(cfName C.CFStringRef) {
		SysexRequest.destination = destination.endpoint         //                        MIDIEndpointRef destination;
		SysexRequest.data = (*C.Byte)(unsafe.Pointer(&data[0])) //                        const Byte *data
		SysexRequest.bytesToSend = (C.UInt32)(len(data))        //                        UInt32 bytesToSend
		//                                                                                Boolean complete;
		//                                                                                Byte reserved[3];
		SysexRequest.completionProc = (C.MIDICompletionProc)(C.getSysexProc()) //         MIDICompletionProc completionProc
		//SysexRequest.completionProc = (C.MIDICompletionProc)(unsafe.Pointer(&MyCallbackFunc)) //         MIDICompletionProc completionProc
		//SysexRequest.completionRefCon = unsafe.Pointer(cfName) //         void *completionRefCon
		SysexRequest.completionRefCon = unsafe.Pointer(&MyCallbackFunc)

	})

	sysexMessage = SysexMessage{SysexRequest, sysexProc}

	return

}

/*


func NewInputPort(client Client, name string, readProc ReadProc) (inputPort InputPort, err error) {
    var port C.MIDIPortRef

    stringToCFString(name, func(cfName C.CFStringRef) {
        osStatus := C.MIDIInputPortCreate(client.client,
            cfName,
            (C.MIDIReadProc)(C.getProc()),       <<
            unsafe.Pointer(uintptr(0)),          <<
            &port)                               <<

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

/*

// for backup
func NewSysexMessage(destination *Destination, data []byte, sysexProc SysexProc) (sysexMessage SysexMessage) {
    var SysexRequest C.MIDISysexSendRequest

    SysexRequest.destination = destination.endpoint
    SysexRequest.data = (*C.Byte)(unsafe.Pointer(&data[0]))
    SysexRequest.bytesToSend = (C.UInt32)(len(data))
    SysexRequest.completionProc = (C.MIDICompletionProc)(C.getSysexProc())

    sysexMessage = SysexMessage{SysexRequest, sysexProc}

    return

}

*/
