package imap

type State struct {
	slots [SlotMax]Slot
	roots [SlotMax]uint32
}
