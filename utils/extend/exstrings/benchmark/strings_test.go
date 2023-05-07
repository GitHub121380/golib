package benchmark

import (
	"strings"
	"testing"

	"github.com/GitHub121380/golib/utils/extend/exstrings"
)

func BenchmarkReplace(b *testing.B) {
	s := "acccbbbaacaabbcbbaaccbaaacaabaccacabbcaacbbccccbbbccaccbcaac"
	for i := 0; i < b.N; i++ {
		exstrings.Replace(s, "cc", "d", -1)
		exstrings.Replace(s, "aa", "d", -1)
		exstrings.Replace(s, "bb", "d", -1)
		exstrings.Replace(s, "ac", "d", -1)
		exstrings.Replace(s, "ca", "d", -1)
		exstrings.Replace(s, "bc", "d", -1)
		exstrings.Replace(s, "ba", "d", -1)
		exstrings.Replace(s, "acc", "d", -1)
		exstrings.Replace(s, "ccb", "d", -1)
		exstrings.Replace(s, "cbb", "d", -1)
		exstrings.Replace(s, "caa", "d", -1)
		exstrings.Replace(s, "bbc", "d", -1)
		exstrings.Replace(s, "aca", "d", -1)
		exstrings.Replace(s, "ccc", "d", -1)
		exstrings.Replace(s, "ab", "d", -1)
		exstrings.Replace(s, "dd", "d", -1)
	}
}

func BenchmarkReplaceToBytes(b *testing.B) {
	s := "acccbbbaacaabbcbbaaccbaaacaabaccacabbcaacbbccccbbbccaccbcaac"
	for i := 0; i < b.N; i++ {
		exstrings.ReplaceToBytes(s, "cc", "d", -1)
		exstrings.ReplaceToBytes(s, "aa", "d", -1)
		exstrings.ReplaceToBytes(s, "bb", "d", -1)
		exstrings.ReplaceToBytes(s, "ac", "d", -1)
		exstrings.ReplaceToBytes(s, "ca", "d", -1)
		exstrings.ReplaceToBytes(s, "bc", "d", -1)
		exstrings.ReplaceToBytes(s, "ba", "d", -1)
		exstrings.ReplaceToBytes(s, "acc", "d", -1)
		exstrings.ReplaceToBytes(s, "ccb", "d", -1)
		exstrings.ReplaceToBytes(s, "cbb", "d", -1)
		exstrings.ReplaceToBytes(s, "caa", "d", -1)
		exstrings.ReplaceToBytes(s, "bbc", "d", -1)
		exstrings.ReplaceToBytes(s, "aca", "d", -1)
		exstrings.ReplaceToBytes(s, "ccc", "d", -1)
		exstrings.ReplaceToBytes(s, "ab", "d", -1)
		exstrings.ReplaceToBytes(s, "dd", "d", -1)
	}
}

func BenchmarkUnsafeReplaceToBytes(b *testing.B) {
	s := "acccbbbaacaabbcbbaaccbaaacaabaccacabbcaacbbccccbbbccaccbcaac"
	for i := 0; i < b.N; i++ {
		exstrings.UnsafeReplaceToBytes(s, "cc", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "aa", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "bb", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "ac", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "ca", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "bc", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "ba", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "acc", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "ccb", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "cbb", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "caa", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "bbc", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "aca", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "ccc", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "ab", "d", -1)
		exstrings.UnsafeReplaceToBytes(s, "dd", "d", -1)
	}
}

func BenchmarkStandardLibraryReplace(b *testing.B) {
	s := "acccbbbaacaabbcbbaaccbaaacaabaccacabbcaacbbccccbbbccaccbcaac"
	for i := 0; i < b.N; i++ {
		strings.Replace(s, "cc", "d", -1)
		strings.Replace(s, "aa", "d", -1)
		strings.Replace(s, "bb", "d", -1)
		strings.Replace(s, "ac", "d", -1)
		strings.Replace(s, "ca", "d", -1)
		strings.Replace(s, "bc", "d", -1)
		strings.Replace(s, "ba", "d", -1)
		strings.Replace(s, "acc", "d", -1)
		strings.Replace(s, "ccb", "d", -1)
		strings.Replace(s, "cbb", "d", -1)
		strings.Replace(s, "caa", "d", -1)
		strings.Replace(s, "bbc", "d", -1)
		strings.Replace(s, "aca", "d", -1)
		strings.Replace(s, "ccc", "d", -1)
		strings.Replace(s, "ab", "d", -1)
		strings.Replace(s, "dd", "d", -1)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exstrings.Repeat("ABC", 100000)
	}
}

func BenchmarkRepeatToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exstrings.RepeatToBytes("ABC", 100000)
	}
}

func BenchmarkStandardLibraryRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Repeat("ABC", 100000)
	}
}

func BenchmarkJoin(b *testing.B) {
	s := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for i := 0; i < b.N; i++ {
		exstrings.Join(s, "-")
	}
}

func BenchmarkJoinToBytes(b *testing.B) {
	s := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for i := 0; i < b.N; i++ {
		exstrings.JoinToBytes(s, "-")
	}
}

func BenchmarkStandardLibraryJoin(b *testing.B) {
	s := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for i := 0; i < b.N; i++ {
		strings.Join(s, "-")
	}
}

func BenchmarkMultiJoinString(b *testing.B) {
	s := []string{"daaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "e", `{"errNo":0,"errStr":"fsd","data":{"105229_0_10":"{\"examInfo\":{\"bindInfo\":{\"bindId\":105229,\"bindType\":0,\"createTime\":1535797800,\"examId\":30188,\"examType\":10,\"operatorName\":\"鲍伟\",\"operatorUid\":18,\"props\":{\"duration\":0,\"maxTryNum\":0,\"passScore\":0},\"relationId\":28418,\"updateTime\":1537023831,\"userKv\":{\"examStatus\":2,\"startTime\":1537023586}},\"examInfo\":\"{\\\"examId\\\":30188,\\\"examType\\\":10,\\\"title\\\":\\\"\\\\u79cb1\\\\u9ad8\\\\u4e09\\\\u7269\\\\u7406\\\\u540c\\\\u6b65(\\\\u5c16\\\\u7aef\\\\u57f9\\\\u4f18\\\\u3001\\\\u5f3a\\\\u5316\\\\u63d0\\\\u5347)\\\\u7b2c2\\\\u8bb2\\\\u5802\\\\u5802\\\\u6d4b\\\",\\\"tidList\\\":{\\\"376085497\\\":{\\\"score\\\":50,\\\"type\\\":2},\\\"375574847\\\":{\\\"score\\\":50,\\\"type\\\":1}},\\\"totalScore\\\":100,\\\"props\\\":[],\\\"userKv\\\":[],\\\"grade\\\":7,\\\"subject\\\":4,\\\"ruleInfo\\\":{\\\"duration\\\":0,\\\"passScore\\\":0,\\\"maxTryNum\\\":0},\\\"extData\\\":[]}\",\"path\":\"c:0-l:0-cpu:0\"},\"questionList\":{}}","105229_0_13":"{\"examInfo\":{},\"questionList\":{}}","105229_0_7":"{\"examInfo\":{\"bindInfo\":{\"bindId\":105229,\"bindType\":0,\"createTime\":1535797998,\"examId\":100724,\"examType\":7,\"operatorName\":\"鲍伟\",\"operatorUid\":18,\"props\":{\"duration\":0,\"maxTryNum\":0,\"passScore\":0},\"relationId\":38712,\"updateTime\":1535797998,\"userKv\":[]},\"examInfo\":\"{\\\"examId\\\":100724,\\\"examType\\\":7,\\\"title\\\":\\\"\\\\u79cb1\\\\u9ad8\\\\u4e09\\\\u7269\\\\u7406\\\\u540c\\\\u6b65(\\\\u5c16\\\\u7aef\\\\u57f9\\\\u4f18\\\\u3001\\\\u5f3a\\\\u5316\\\\u63d0\\\\u5347)\\\\u7b2c2\\\\u8bb2\\\\u8bfe\\\\u540e\\\\u4f5c\\\\u4e1a\\\",\\\"tidList\\\":{\\\"375479289\\\":{\\\"score\\\":0,\\\"type\\\":2},\\\"375574847\\\":{\\\"score\\\":0,\\\"type\\\":1}},\\\"totalScore\\\":0,\\\"props\\\":[],\\\"userKv\\\":[],\\\"grade\\\":7,\\\"subject\\\":4,\\\"ruleInfo\\\":{\\\"duration\\\":0,\\\"passScore\\\":0,\\\"maxTryNum\\\":0},\\\"extData\\\":[]}\",\"path\":\"c:0-l:0-cpu:0\"},\"questionList\":{}}"}}`, "g", "hdas", "iadsa", "j", "k", "l", "mdsa", "nad", "oadadasdasdsa"}
	for i := 0; i < b.N; i++ {
		exstrings.MultiJoinString(s...)
	}
}

// func BenchmarkMultiJoinStringByBufferPool(b *testing.B) {
// 	s := []string{"daaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "e", `{"errNo":0,"errStr":"fsd","data":{"105229_0_10":"{\"examInfo\":{\"bindInfo\":{\"bindId\":105229,\"bindType\":0,\"createTime\":1535797800,\"examId\":30188,\"examType\":10,\"operatorName\":\"鲍伟\",\"operatorUid\":18,\"props\":{\"duration\":0,\"maxTryNum\":0,\"passScore\":0},\"relationId\":28418,\"updateTime\":1537023831,\"userKv\":{\"examStatus\":2,\"startTime\":1537023586}},\"examInfo\":\"{\\\"examId\\\":30188,\\\"examType\\\":10,\\\"title\\\":\\\"\\\\u79cb1\\\\u9ad8\\\\u4e09\\\\u7269\\\\u7406\\\\u540c\\\\u6b65(\\\\u5c16\\\\u7aef\\\\u57f9\\\\u4f18\\\\u3001\\\\u5f3a\\\\u5316\\\\u63d0\\\\u5347)\\\\u7b2c2\\\\u8bb2\\\\u5802\\\\u5802\\\\u6d4b\\\",\\\"tidList\\\":{\\\"376085497\\\":{\\\"score\\\":50,\\\"type\\\":2},\\\"375574847\\\":{\\\"score\\\":50,\\\"type\\\":1}},\\\"totalScore\\\":100,\\\"props\\\":[],\\\"userKv\\\":[],\\\"grade\\\":7,\\\"subject\\\":4,\\\"ruleInfo\\\":{\\\"duration\\\":0,\\\"passScore\\\":0,\\\"maxTryNum\\\":0},\\\"extData\\\":[]}\",\"path\":\"c:0-l:0-cpu:0\"},\"questionList\":{}}","105229_0_13":"{\"examInfo\":{},\"questionList\":{}}","105229_0_7":"{\"examInfo\":{\"bindInfo\":{\"bindId\":105229,\"bindType\":0,\"createTime\":1535797998,\"examId\":100724,\"examType\":7,\"operatorName\":\"鲍伟\",\"operatorUid\":18,\"props\":{\"duration\":0,\"maxTryNum\":0,\"passScore\":0},\"relationId\":38712,\"updateTime\":1535797998,\"userKv\":[]},\"examInfo\":\"{\\\"examId\\\":100724,\\\"examType\\\":7,\\\"title\\\":\\\"\\\\u79cb1\\\\u9ad8\\\\u4e09\\\\u7269\\\\u7406\\\\u540c\\\\u6b65(\\\\u5c16\\\\u7aef\\\\u57f9\\\\u4f18\\\\u3001\\\\u5f3a\\\\u5316\\\\u63d0\\\\u5347)\\\\u7b2c2\\\\u8bb2\\\\u8bfe\\\\u540e\\\\u4f5c\\\\u4e1a\\\",\\\"tidList\\\":{\\\"375479289\\\":{\\\"score\\\":0,\\\"type\\\":2},\\\"375574847\\\":{\\\"score\\\":0,\\\"type\\\":1}},\\\"totalScore\\\":0,\\\"props\\\":[],\\\"userKv\\\":[],\\\"grade\\\":7,\\\"subject\\\":4,\\\"ruleInfo\\\":{\\\"duration\\\":0,\\\"passScore\\\":0,\\\"maxTryNum\\\":0},\\\"extData\\\":[]}\",\"path\":\"c:0-l:0-cpu:0\"},\"questionList\":{}}"}}`, "g", "hdas", "iadsa", "j", "k", "l", "mdsa", "nad", "oadadasdasdsa"}

//		for i := 0; i < b.N; i++ {
//			exstrings.MultiJoinStringByBufferPool(s...)
//		}
//	}
func BenchmarkJI_MultiJoinString(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数
	s := []string{"daaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "e", `{"errNo":0,"errStr":"fsd","data":{"105229_0_10":"{\"examInfo\":{\"bindInfo\":{\"bindId\":105229,\"bindType\":0,\"createTime\":1535797800,\"examId\":30188,\"examType\":10,\"operatorName\":\"鲍伟\",\"operatorUid\":18,\"props\":{\"duration\":0,\"maxTryNum\":0,\"passScore\":0},\"relationId\":28418,\"updateTime\":1537023831,\"userKv\":{\"examStatus\":2,\"startTime\":1537023586}},\"examInfo\":\"{\\\"examId\\\":30188,\\\"examType\\\":10,\\\"title\\\":\\\"\\\\u79cb1\\\\u9ad8\\\\u4e09\\\\u7269\\\\u7406\\\\u540c\\\\u6b65(\\\\u5c16\\\\u7aef\\\\u57f9\\\\u4f18\\\\u3001\\\\u5f3a\\\\u5316\\\\u63d0\\\\u5347)\\\\u7b2c2\\\\u8bb2\\\\u5802\\\\u5802\\\\u6d4b\\\",\\\"tidList\\\":{\\\"376085497\\\":{\\\"score\\\":50,\\\"type\\\":2},\\\"375574847\\\":{\\\"score\\\":50,\\\"type\\\":1}},\\\"totalScore\\\":100,\\\"props\\\":[],\\\"userKv\\\":[],\\\"grade\\\":7,\\\"subject\\\":4,\\\"ruleInfo\\\":{\\\"duration\\\":0,\\\"passScore\\\":0,\\\"maxTryNum\\\":0},\\\"extData\\\":[]}\",\"path\":\"c:0-l:0-cpu:0\"},\"questionList\":{}}","105229_0_13":"{\"examInfo\":{},\"questionList\":{}}","105229_0_7":"{\"examInfo\":{\"bindInfo\":{\"bindId\":105229,\"bindType\":0,\"createTime\":1535797998,\"examId\":100724,\"examType\":7,\"operatorName\":\"鲍伟\",\"operatorUid\":18,\"props\":{\"duration\":0,\"maxTryNum\":0,\"passScore\":0},\"relationId\":38712,\"updateTime\":1535797998,\"userKv\":[]},\"examInfo\":\"{\\\"examId\\\":100724,\\\"examType\\\":7,\\\"title\\\":\\\"\\\\u79cb1\\\\u9ad8\\\\u4e09\\\\u7269\\\\u7406\\\\u540c\\\\u6b65(\\\\u5c16\\\\u7aef\\\\u57f9\\\\u4f18\\\\u3001\\\\u5f3a\\\\u5316\\\\u63d0\\\\u5347)\\\\u7b2c2\\\\u8bb2\\\\u8bfe\\\\u540e\\\\u4f5c\\\\u4e1a\\\",\\\"tidList\\\":{\\\"375479289\\\":{\\\"score\\\":0,\\\"type\\\":2},\\\"375574847\\\":{\\\"score\\\":0,\\\"type\\\":1}},\\\"totalScore\\\":0,\\\"props\\\":[],\\\"userKv\\\":[],\\\"grade\\\":7,\\\"subject\\\":4,\\\"ruleInfo\\\":{\\\"duration\\\":0,\\\"passScore\\\":0,\\\"maxTryNum\\\":0},\\\"extData\\\":[]}\",\"path\":\"c:0-l:0-cpu:0\"},\"questionList\":{}}"}}`, "g", "hdas", "iadsa", "j", "k", "l", "mdsa", "nad", "oadadasdasdsa"}
	b.StartTimer() //重新开始时间
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			exstrings.MultiJoinString(s...)
		}
	})
}

// func BenchmarkJI_MultiJoinStringByBufferPool(b *testing.B) {
// 	b.StopTimer() //调用该函数停止压力测试的时间计数
// 	s := []string{"daaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "e", `{"errNo":0,"errStr":"fsd","data":{"105229_0_10":"{\"examInfo\":{\"bindInfo\":{\"bindId\":105229,\"bindType\":0,\"createTime\":1535797800,\"examId\":30188,\"examType\":10,\"operatorName\":\"鲍伟\",\"operatorUid\":18,\"props\":{\"duration\":0,\"maxTryNum\":0,\"passScore\":0},\"relationId\":28418,\"updateTime\":1537023831,\"userKv\":{\"examStatus\":2,\"startTime\":1537023586}},\"examInfo\":\"{\\\"examId\\\":30188,\\\"examType\\\":10,\\\"title\\\":\\\"\\\\u79cb1\\\\u9ad8\\\\u4e09\\\\u7269\\\\u7406\\\\u540c\\\\u6b65(\\\\u5c16\\\\u7aef\\\\u57f9\\\\u4f18\\\\u3001\\\\u5f3a\\\\u5316\\\\u63d0\\\\u5347)\\\\u7b2c2\\\\u8bb2\\\\u5802\\\\u5802\\\\u6d4b\\\",\\\"tidList\\\":{\\\"376085497\\\":{\\\"score\\\":50,\\\"type\\\":2},\\\"375574847\\\":{\\\"score\\\":50,\\\"type\\\":1}},\\\"totalScore\\\":100,\\\"props\\\":[],\\\"userKv\\\":[],\\\"grade\\\":7,\\\"subject\\\":4,\\\"ruleInfo\\\":{\\\"duration\\\":0,\\\"passScore\\\":0,\\\"maxTryNum\\\":0},\\\"extData\\\":[]}\",\"path\":\"c:0-l:0-cpu:0\"},\"questionList\":{}}","105229_0_13":"{\"examInfo\":{},\"questionList\":{}}","105229_0_7":"{\"examInfo\":{\"bindInfo\":{\"bindId\":105229,\"bindType\":0,\"createTime\":1535797998,\"examId\":100724,\"examType\":7,\"operatorName\":\"鲍伟\",\"operatorUid\":18,\"props\":{\"duration\":0,\"maxTryNum\":0,\"passScore\":0},\"relationId\":38712,\"updateTime\":1535797998,\"userKv\":[]},\"examInfo\":\"{\\\"examId\\\":100724,\\\"examType\\\":7,\\\"title\\\":\\\"\\\\u79cb1\\\\u9ad8\\\\u4e09\\\\u7269\\\\u7406\\\\u540c\\\\u6b65(\\\\u5c16\\\\u7aef\\\\u57f9\\\\u4f18\\\\u3001\\\\u5f3a\\\\u5316\\\\u63d0\\\\u5347)\\\\u7b2c2\\\\u8bb2\\\\u8bfe\\\\u540e\\\\u4f5c\\\\u4e1a\\\",\\\"tidList\\\":{\\\"375479289\\\":{\\\"score\\\":0,\\\"type\\\":2},\\\"375574847\\\":{\\\"score\\\":0,\\\"type\\\":1}},\\\"totalScore\\\":0,\\\"props\\\":[],\\\"userKv\\\":[],\\\"grade\\\":7,\\\"subject\\\":4,\\\"ruleInfo\\\":{\\\"duration\\\":0,\\\"passScore\\\":0,\\\"maxTryNum\\\":0},\\\"extData\\\":[]}\",\"path\":\"c:0-l:0-cpu:0\"},\"questionList\":{}}"}}`, "g", "hdas", "iadsa", "j", "k", "l", "mdsa", "nad", "oadadasdasdsa"}
// 	b.StartTimer() //重新开始时间
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			exstrings.MultiJoinStringByBufferPool(s...)
// 		}
// 	})
// }
