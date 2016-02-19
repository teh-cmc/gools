package zerocopy

import "testing"

// -----------------------------------------------------------------------------

func TestZerocopy_StringToByteSlice(t *testing.T) {
	var str string
	var expected []byte

	str = ""
	expected = []byte(str)
	if ret := StringToByteSlice(str); string(ret) != string(expected) {
		t.Errorf("expected %v, got %v", expected, ret)
	}

	str = "Ṱ̺̺̕o͞ ̷i̲̬͇̪͙n̝̗͕v̟̜̘̦͟o̶̙̰̠kè͚̮̺̪̹̱̤ ̖t̝͕̳̣̻̪͞h̼͓̲̦̳̘̲e͇̣̰̦̬͎ ̢̼̻̱̘h͚͎͙̜̣̲ͅi̦̲̣̰̤v̻͍e̺̭̳̪̰-m̢iͅn̖̺̞̲̯̰d̵̼̟͙̩̼̘̳ ̞̥̱̳̭r̛̗̘e͙p͠r̼̞̻̭̗e̺̠̣͟s̘͇̳͍̝͉e͉̥̯̞̲͚̬͜ǹ̬͎͎̟̖͇̤t͍̬̤͓̼̭͘ͅi̪̱n͠g̴͉ ͏͉ͅc̬̟h͡a̫̻̯͘o̫̟̖͍̙̝͉s̗̦̲.̨̹͈̣"
	expected = []byte(str)
	if ret := StringToByteSlice(str); string(ret) != string(expected) {
		t.Errorf("expected %v, got %v", expected, ret)
	}
}
