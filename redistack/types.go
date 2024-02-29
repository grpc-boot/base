package redistack

import (
	"github.com/redis/rueidis"
)

type FieldValue map[string]string

type InfoStream struct {
	Length             int64
	Groups             int64
	MaxDeletedEntryId  string
	RecordFirstEntryId string
	LastGeneratedId    string
	FirstEntry         rueidis.XRangeEntry
	LastEntry          rueidis.XRangeEntry
	EntriesAdded       int64
	RadixTreeKeys      int64
	RadixTreeNodes     int64
}

func ToInfoStream(msgMap map[string]rueidis.RedisMessage) InfoStream {
	length := msgMap["length"]
	groups := msgMap["groups"]
	maxDeletedEntryId := msgMap["max-deleted-entry-id"]
	recordFirstEntryId := msgMap["recorded-first-entry-id"]
	lastGeneratedId := msgMap["last-generated-id"]
	firstEntry := msgMap["first-entry"]
	lastEntry := msgMap["last-entry"]
	entriesAdded := msgMap["entries-added"]
	radixTreeKeys := msgMap["radix-tree-keys"]
	radixTreeNodes := msgMap["radix-tree-nodes"]

	info := InfoStream{}
	info.Groups, _ = groups.AsInt64()
	info.Length, _ = length.AsInt64()
	info.MaxDeletedEntryId, _ = maxDeletedEntryId.ToString()
	info.RecordFirstEntryId, _ = recordFirstEntryId.ToString()
	info.LastGeneratedId, _ = lastGeneratedId.ToString()
	info.FirstEntry, _ = firstEntry.AsXRangeEntry()
	info.LastEntry, _ = lastEntry.AsXRangeEntry()
	info.EntriesAdded, _ = entriesAdded.AsInt64()
	info.RadixTreeKeys, _ = radixTreeKeys.AsInt64()
	info.RadixTreeNodes, _ = radixTreeNodes.AsInt64()

	return info
}
