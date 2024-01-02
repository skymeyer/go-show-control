package artnet

//go:generate stringer -linecomment -type OemCode
//go:generate stringer -linecomment -type ESTAManCode
//go:generate stringer -type OpCode

type OemCode uint16 // TODO: import OEM codes from source

const (
	OemDMXHub OemCode = 0x0000 // Manufacturer: Artistic Licence Engineering Ltd ProductName: Dmx Hub   NumDmxIn: 4 NumDmxOut: 4 DmxPortPhysical: y RdmSupported: n SupportEmail: support@ArtisticLicence.com SupportName: Wayne Howell CoWeb: http://www.ArtisticLicence.com
	OemGlobal OemCode = 0xffff // Manufacturer: Artistic Licence Engineering Ltd ProductName: OemGlobal NumDmxIn: 0 NumDmxOut: 0 DmxPortPhysical: n RdmSupported: n SupportEmail: support@ArtisticLicence.com SupportName: Wayne Howell CoWeb: http://www.ArtisticLicence.com
)

type ESTAManCode uint16 // TODO: import ESTA Manufacturer codes from source

const (
	ESTAManUnknown ESTAManCode = 0x0000
	ESTAManGlobal  ESTAManCode = 0xffff
)

type OpCode uint16

const (
	OpPoll      OpCode = 0x2000 // This is an ArtPoll packet, no other data is contained in this UDP packet.
	OpPollReply OpCode = 0x2100 // This is an ArtPollReply Packet. It contains device status information.
	OpDiagData  OpCode = 0x2300 // Diagnostics and data logging packet.
	OpCommand   OpCode = 0x2400 // This is an ArtCommand packet. It is used to send text based parameter commands.

	OpDataRequest OpCode = 0x2700 // This is an ArtDataRequest packet. It is used to request data such as products URLs
	OpDataReply   OpCode = 0x2800 // This is an ArtDataReply packet. It is used to reply to ArtDataRequest packets.

	OpDmx  OpCode = 0x5000 // This is an ArtDmx data packet. It contains zero start code DMX512 information for a single Universe.
	OpNzs  OpCode = 0x5100 // This is an ArtNzs data packet. It contains non-zero start code (except RDM) DMX512 information for a single Universe.
	OpSync OpCode = 0x5200 // This is an ArtSync data packet. It is used to force synchronous transfer of ArtDmx packets to a node’s output.

	OpAddress    OpCode = 0x6000 // This is an ArtAddress packet. It contains remote programming information for a Node.
	OpInput      OpCode = 0x7000 // This is an ArtInput packet. It contains enable – disable data for DMX inputs.
	OpTodRequest OpCode = 0x8000 // This is an ArtTodRequest packet. It is used to request a Table of Devices (ToD) for RDM discovery.
	OpTodData    OpCode = 0x8100 // This is an ArtTodData packet. It is used to send a Table of Devices (ToD) for RDM discovery.
	OpTodControl OpCode = 0x8200 // This is an ArtTodControl packet. It is used to send RDM discovery control messages.
	OpRdm        OpCode = 0x8300 // This is an ArtRdm packet. It is used to send all non discovery RDM messages.

	OpRdmSub       OpCode = 0x8400 // This is an ArtRdmSub packet. It is used to send compressed, RDM Sub-Device data.
	OpVideoSetup   OpCode = 0xa010 // This is an ArtVideoSetup packet. It contains video screen setup information for nodes that implement the extended video features.
	OpVideoPalette OpCode = 0xa020 // This is an ArtVideoPalette packet. It contains colour palette setup information for nodes that implement the extended video features.
	OpVideoData    OpCode = 0xa040 // This is an ArtVideoData packet. It contains display data for nodes that implement the extended video features.

	OpMacMaster OpCode = 0xf000 // This packet is deprecated.
	OpMacSlave  OpCode = 0xf100 // This packet is deprecated.

	OpFirmwareMaster   OpCode = 0xf200 // This is an ArtFirmwareMaster packet. It is used to upload new firmware or firmware extensions to the Node.
	OpFirmwareReply    OpCode = 0xf300 // This is an ArtFirmwareReply packet. It is returned by the node to acknowledge receipt of an ArtFirmwareMaster packet or ArtFileTnMaster packet.
	OpFileTnMaster     OpCode = 0xf400 // Uploads user file to node.
	OpFileFnMaster     OpCode = 0xf500 // Downloads user file from node.
	OpFileFnReply      OpCode = 0xf600 // Server to Node acknowledge for download packets.
	OpIpProg           OpCode = 0xf800 // This is an ArtIpProg packet. It is used to reprogramme the IP address and Mask of the Node.
	OpIpProgReply      OpCode = 0xf900 // This is an ArtIpProgReply packet. It is returned by the node to acknowledge receipt of an ArtIpProg packet.
	OpMedia            OpCode = 0x9000 // This is an ArtMedia packet. It is Unicast by a Media Server and acted upon by a Controller.
	OpMediaPatch       OpCode = 0x9100 // This is an ArtMediaPatch packet. It is Unicast by a Controller and acted upon by a Media Server.
	OpMediaControl     OpCode = 0x9200 // This is an ArtMediaControl packet. It is Unicast by a Controller and acted upon by a Media Server.
	OpMediaContrlReply OpCode = 0x9300 // This is an ArtMediaControlReply packet. It is Unicast by a Media Server and acted upon by a Controller.

	OpTimeCode       OpCode = 0x9700 //This is an ArtTimeCode packet. It is used to transport time code over the network.
	OpTimeSync       OpCode = 0x9800 // Used to synchronise real time date and clock
	OpTrigger        OpCode = 0x9900 // Used to send trigger macros
	OpDirectory      OpCode = 0x9a00 // Requests a node's file list
	OpDirectoryReply OpCode = 0x9b00 // Replies to OpDirectory with file list
)

type ReportCode uint16

const (
	RcDebug        ReportCode = 0x0000 // Booted in debug mode (Only used in development)
	RcPowerOk      ReportCode = 0x0001 // Power On Tests successful
	RcPowerFail    ReportCode = 0x0002 // Hardware tests failed at Power On
	RcSocketWr1    ReportCode = 0x0003 // Last UDP from Node failed due to truncated length,  Most likely caused by a collision.
	RcParseFail    ReportCode = 0x0004 // Unable to identify last UDP transmission. Check OpCode and packet length.
	RcUdpFail      ReportCode = 0x0005 // Unable to open Udp Socket in last transmission attempt
	RcShNameOk     ReportCode = 0x0006 // Confirms that Port Name programming via ArtAddress, was successful.
	RcLoNameOk     ReportCode = 0x0007 // Confirms that Long Name programming via ArtAddress, was successful.
	RcDmxError     ReportCode = 0x0008 // DMX512 receive errors detected.
	RcDmxUdpFull   ReportCode = 0x0009 // Ran out of internal DMX transmit buffers.
	RcDmxRxFull    ReportCode = 0x000a // Ran out of internal DMX Rx buffers.
	RcSwitchErr    ReportCode = 0x000b // Rx Universe switches conflict.
	RcConfigErr    ReportCode = 0x000c // Product configuration does not match firmware.
	RcDmxShort     ReportCode = 0x000d // DMX output short detected. See GoodOutput field.
	RcFirmwareFail ReportCode = 0x000e //Last attempt to upload new firmware failed.
	RcUserFail     ReportCode = 0x000f // User changed switch settings when address locked by remote programming. User changes ignored.
	RcFactoryRes   ReportCode = 0x0010 //Factory reset has occurred.

)

type StyleCode uint8

const (
	StNode       StyleCode = 0x00 // A DMX to / from Art-Net device
	StController StyleCode = 0x01 // A lighting console.
	StMedia      StyleCode = 0x02 // A Media Server.
	StRoute      StyleCode = 0x03 // A network routing device.
	StBackup     StyleCode = 0x04 // A backup device.
	StConfig     StyleCode = 0x05 // A configuration or diagnostic tool.
	StVisual     StyleCode = 0x06 //A visualiser.
)

// Diagnostics Priority Code
type DpCode uint8

const (
	DpLow      DpCode = 0x10 // Low priority message.
	DpMed      DpCode = 0x40 // Medium priority message.
	DpHigh     DpCode = 0x80 // High priority message.
	DpCritical DpCode = 0xe0 // Critical priority message.
	DpVolatile DpCode = 0xf0 // Volatile message.
)

type TalkToMe uint8

func NewTalkToMe() TalkToMe {
	return TalkToMe(0)
}

func (t TalkToMe) With(f TalkToMeFlag) TalkToMe {
	return TalkToMe(uint8(t) | uint8(f))
}

func (t TalkToMe) Without(f TalkToMeFlag) TalkToMe {
	return TalkToMe(uint8(t) &^ uint8(f))
}

func (t TalkToMe) Has(f TalkToMeFlag) bool {
	return uint8(t)&uint8(f) != 0
}

type TalkToMeFlag uint8

const (
	TTMReplyOnChange      TalkToMeFlag = 2 << iota // Send ArtPollReply whenever Node conditions change.
	TTMDiagnostics                                 // Send me diagnostics messages.
	TTMDiagnosticsUnicast                          // Diagnostics messages are unicast if enabled.
	TTMDisableVLC                                  // Disable VLC transmission.
	TTMEnableTargetedMode                          // Enable Targeted Mode.
)
