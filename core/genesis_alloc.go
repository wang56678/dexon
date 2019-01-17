// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

// Constants containing the genesis allocation of built-in genesis blocks.
// Their content is an RLP-encoded list of (address, balance) tuples.
// Use mkalloc.go to create/update them.

// nolint: misspell
const mainnetAllocData = testnetAllocData

const testnetAllocData = "\xf9\x03>\xf8\xb8\x94L\x05N\xaeq\x18\xe4\xe8Mf\x0f\u007f\u076c\x98W>-;~\xf8\xa1\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\x1f\x89\u06c1\xf3\xc7\t9S_\xf0<f{\xc3\xd4T)\xe7>?H:Wxo\fB\xc2\xcei\xbe'\xe8\xe7X\xbb\x1a\xaa\xfd\x17J:k\u01ac\x88.\xe6)\u0703\x02\xd8\x19)K\xb2%#\u030e\xbb\xe5\xf8D\x91DEXON Test Node 3\x90dexon3@dexon.org\x8eTaipei, Taiwan\x91https://dexon.org\xf8\xb8\x94e=\rg\xff\xe9\xfd\xb3\x1c\xa7\xaf\xeaI\xc62~\xed\x11n\xa7\xf8\xa1\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xacb\xf0\n\xa3\xc3\xebZ\x1cAh\xdem/\x9a\x05\xf2\xc7\x1a\n\bk\u0720\xfd\xaa\xff&\xde\x18@=(\xfaS\x1e\"-ku\x96n\xe4y\u5a35\xd5\xcdB?\xb3+\xe8qQC\xf3\nJ\xd5\xd77\xa0\xf8D\x91DEXON Test Node 0\x90dexon0@dexon.org\x8eTaipei, Taiwan\x91https://dexon.org\xf8\xb8\x94\x98FK\xa2#p\x9bq\x8c\xef\x05\u067cA\xc9\xf4\xa8\xbcH`\xf8\xa1\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xa1P\x17~\u0372\u05dd\xdf\u007f\xe2-X\x91\x8b\n\x98\xa8\x8e\xd2\u0308\x8b\xa3W\x90\xfauFY\x80U\\p\xc9F\xefs=\xfb'\x86N2\xe6\xe85b(\xb6^@\x84,\xb1mE\xbd\xe1G8\xe5\xb6T\xf8D\x91DEXON Test Node 1\x90dexon1@dexon.org\x8eTaipei, Taiwan\x91https://dexon.org\xf8\xb8\x94\xa0\x8f\xb1\xa7e$\x1d\x1cd\xbf\x80\xae\xca\u0414\u05afz%L\xf8\xa1\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xeb\x18\x16^\xeb\x02\x11\xa6.\x96P[\xd5\xd3i\xc7-\u03a0\xfc\x897U\x8a\x01\x10\u007f\xddPA\xc4\u0205o\x0f \x98\xb0K\xb8\xc4X\x13\x0e\x11\aL\x1et\x1f\xd4\xf4\xad.\x88\u04a1\xf2\xb0WG\xf7\x04\x06\xf8D\x91DEXON Test Node 2\x90dexon2@dexon.org\x8eTaipei, Taiwan\x91https://dexon.org\ua53f\x8cH\xa6 \xba\xccF\x90\u007f\x9b\x89s-%\xe4z-|\xf7\u050b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x80\x80\x80\u0100\x80\x80\x80\xea\x94\xe0\xf8Y4\x03S\x85F\x93\xf57\xea3q\x06\xa3>\x9f\xea\xb0\u050b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x80\x80\x80\u0100\x80\x80\x80"

const taipeiAllocData = "\xf9\x18\x1c\xf8\xbe\x93RqT=FG\xab8\x8b\xfc\x9c\x8f\xda8\x9a+\xfbu\x83\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xc8\xf5\vgh\u02ab0\x8c\x1fhp\xc4'bM\xcf*6i\x96\x90\xe4o\xca\x10Ik-\x18\x8f\x81[\\z\xb5\xf5e\xdd\xfd\x89\xfbH\xf8\x92\x99\xe7K5\x95\x98\ag\xbf\xa6V\xf6\x8bN0O(C\xf5\xf8K\x91Node - us-west1-1\x9btaipei-us-west1-1@dexon.org\x8aus-west1-1\x91https://dexon.org\xf8\xbf\x94\x1b\x03Ja\x83\x93\x96?!\xd3\xd7HOh\x164\u05d8\xec\x0f\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04#\xd4\fYTw\xdf(\xc8E@\u0742$n\xfa\x9f\xa6d\xab9?07\x91Y\x00G\xe8)~D\x86\xbb\xe1\x81\xe0\xdd\x12\xfc AG\x91n\x94V\x8cO\xb6\xb4\x91\xae\xfb\r\u0564[\xe6E\xd7\xd2@<\xf8K\x91Node - us-east1-1\x9btaipei-us-east1-1@dexon.org\x8aus-east1-1\x91https://dexon.org\xf8\xec\x94\x1c\x9b\u007f\xeb\xefdc\\\xe8\xd0\x17\x02\x8e\x82+CdJ\v\xff\xf8\u054b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xa8\x1d\x931 \x9a_\x8c\xebY5\x91\x19\x1b\x1c\xfd`\x96n\xc8\u0583\x06\xf8\xa0\xeb\xacE\x8a\x112U\xceE\x90\x1a`\xb7\x18\xff\x9e\xb3\xcb\x16\xb8\x82b\xaaf\x1f\xd6\xed\x1eu\xfb- \xba+\u0266\xaf\xb6\r\xf8x\xa0Node - northamerica-northeast1-0\xaataipei-northamerica-northeast1-0@dexon.org\x99northamerica-northeast1-0\x91https://dexon.org\xf8\u02d4#C\xb4\x947S\x06\xf7\xf9\xa2\x82|SEG\xdf\u05faw~\xf8\xb4\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04@\xe6\x10\u0474\x06_\xbb\xe6\xf63\xae$\u01bb\x86\x14\xe9\x10W\xc7\xda\x1d\xb5\x8ew\xa9t\x1elw\xf4\v\xab\xa7d\u068c\xc4\u0768\x16&\x83sE\xceK\xc0F\x1ea?\xab\xbf\xe1o\xf2\x80\xe4uW\x8d\xee\xf8W\x95Node - europe-west4-1\x9ftaipei-europe-west4-1@dexon.org\x8eeurope-west4-1\x91https://dexon.org\xf8\xbf\x94'\xba\xa0\x95\xbd\x1d\x1d\u04f5\x1do\x99\x15\xdb \x03\xd1\xd5&\\\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xbb\xac\xf8\xd4%:\x18\x1e'\xd3!W\x19\x99\x90\x99A\xbd\xd0e0Q\xb2L\x1c\xd1\xf0\xbe\xfa\x86p\xac\x91\x8d\u00ccd\xf4&\xdc\xde\u37b9xQ\x9d\x16\xd8\x0f\x84\xaa\x00\x05\xf1i\x1c\x96\xff\xee\xd3-\x0e(\xf8K\x91Node - us-east4-0\x9btaipei-us-east4-0@dexon.org\x8aus-east4-0\x91https://dexon.org\xf8\xbf\x94(\x88\u052eUi\u0168\xeb\xae\xc6\xc3!Ycy\x8c8\xe5\x9b\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04+\xab!3\x86\x99U\x9a\xd7\"\v.e\xe7N\xb8B\xf7\x13\v\xaa\xc0\xa7\xeab\xa9\u007f(\x18\"l\n\x8e\xf04\xa1\xe6Z\xfd\r\x90\u07a17G\a\u05bd\r\x0fR\xb12\x15.\x1f CS\xbe\xb1v\x9b\xbb\xf8K\x91Node - us-east4-1\x9btaipei-us-east4-1@dexon.org\x8aus-east4-1\x91https://dexon.org\xf8\u0154.SN\x10|3\xb8\xfb.\x96x\xae\xb1\x86F\xda\x1d\xa0\x1b\x8d\xf8\xae\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\x9b\x8d\xc3'\xa1Q\xe0\u068e\u024d\xef\a\xdc\x05;-1\x86\x19/\x1d\xa9\xfaY\u0345p\xc7x\xf6e\x98\xa7L\xf9\x91\x83\x9c D\xb5\xeb\x14\xa2e\x0e\x86\x19n\xc5\u04ec\xfa\x1d7\xfb\x0e\xd2\xd7\xf3i\\\x04\xf8Q\x93Node - asia-east1-2\x9dtaipei-asia-east1-2@dexon.org\x8casia-east1-2\x91https://dexon.org\xf8\xbf\x94:\xb3o\u007f\t\xa4v\u007f\xa7\x97\xf8\xe0\x96\x0e{$~\x0e\xfe\x02\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xa8*\x8c\x16[\x96}\"\xeb\x83w\x93\x15\x04X\xe7\xf7\vhIj\xff\xe9\xa4\xf0,\x99\x04b\x81\x93'\x86\xab\xfb\xfd\xa2u\xc7\r=$\x94L@\x92T\xe0Jrc\xe2\xdc\xf9\x85\u07cc\xae\"ir;\xfd6\xf8K\x91Node - us-west1-2\x9btaipei-us-west1-2@dexon.org\x8aus-west1-2\x91https://dexon.org\xf8\u02d4I|\xef8/\xbd\x90\u007f_\x90\x89!\xee\xe8\u07f021<\xa7\xf8\xb4\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04VVw\x9e)\x99\xf1;\xb6\x86\x1c\u048f\xbfo>\x88\xef\x91\xfb\b\xc6\u0de9j\x95\x9f\xc0h\x0e\xdc7\"\x16\xdf!;C\x91M\x05a\v\xbb\x06\x17\x1fU\xf83\x1a.\\\xff'e@L\xa4\xd5=\xec\xaf\xf8W\x95Node - europe-west4-2\x9ftaipei-europe-west4-2@dexon.org\x8eeurope-west4-2\x91https://dexon.org\xf8\xbf\x94PY\x8e}\x86\x8d^;\x8c;\u1ff0Pt'}\xf05d\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\x9e\x00\xde\xf0\x98\x1fN\xda\u04eci0\x88\x16\u0095\r_\x15n [\xed\xd3g\xf4\xf4\xccw\xb3\xd3+S\x0f\u0214\v\x97;\xf0$s\xf4\xc1\a'^\xc9\xf3\xa7\xf1\n\x19\xfa\xa5\xcf7\xd0hj\x16\u007f\x9b\xa7\xf8K\x91Node - us-west1-0\x9btaipei-us-west1-0@dexon.org\x8aus-west1-0\x91https://dexon.org\xf8\xbf\x94W\x82\u483eN$Z\xb7-c\u03a5J\x10Ag\u0392\xc0\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04xk%Y\x99\fd\u0451\xacqg\x04\x10.\x16\x95q\x90\x81a\xa2W)^9#\x12\x1b\x15ZC\xae\xaa\x9d\x97\r\xf7)c\x93\xc7\xcddI,,\x1d\x90|yN$`\xfdL8\xab\x05\xb7M\u00ee\xf1\xf8K\x91Node - us-east1-2\x9btaipei-us-east1-2@dexon.org\x8aus-east1-2\x91https://dexon.org\xf8\xbf\x94]\xc3\x04\u03c2\xa1\xf01n\x955\x83\x9f]x\x01\xeft\x02\x1e\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04q\xa2:v\xee\xa5'\xc5\x00\xe0\xbcY\xbbQ3\xa9\xfc\xa7\xff\x16(\xa9 \xee:d\x9aR{C\xad\x05dU\xc5\u04aa\xf7\u03e3<\xef\x18|\x89X1\xfb\xd1=\xa9\x1c\xf4}1\xb5\xe2\x18\xd7-$3a\x0f\xf8K\x91Node - us-west2-1\x9btaipei-us-west2-1@dexon.org\x8aus-west2-1\x91https://dexon.org\xf8\u0514_*(\xa5\xc3\x14\x8ar\xe1j\fD\xa7\xb0\xce\xc1/\v\xa2\x17\xf8\xbd\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04Js\x93\xd7\x10BG\xf2\xe8R\x90d\xac-\xcbi\x91\xec$\xc7\x14\xdb\xfawJ\u01eb\x00\xd9\u0552\x81f\xac1\x89u\xd1\u05bcQ\xc5\x13[\xd2i\xcb\xc5\xee\xfc\xcahY\xba0\xa4\x97\x95j\xb5\xf6\x9d\x14\xba\xf8`\x98Node - asia-northeast1-2\xa2taipei-asia-northeast1-2@dexon.org\x91asia-northeast1-2\x91https://dexon.org\xf8\u0514j=d>|N\xdf3[\f\x89e\x99\bT:P\xee\xbbJ\xf8\xbd\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04d(\xc3&\x15o\x1d\xa9\xd6\xe4\x06\xdc\xc9\xfc\xd5\xeb\xd5(J\x01\xbcB\xbco\nf)\x04\x00\u01ad\x9a\x01\x803\xaa\xf5h/9\x9a\xc0\xb5\xed\xa4C\xf4\xdf\xfd\x17\xc6\x1fH\x15\x8eq\xf6\xf9\xbd\u015e\xa2-\x8b\xf8`\x98Node - asia-northeast1-1\xa2taipei-asia-northeast1-1@dexon.org\x91asia-northeast1-1\x91https://dexon.org\xf8\u0214\x89\xba\u0241K\xd8J\x85\xe5\x87\xda\v\xbf\x88l=\x11\x98\"+\xf8\xb1\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\x81\x8eK\xfc\xd6O\u00f3\xdf*n\t\xa4\xf4\x90\xb9\x86E\xff\xc0\xf6E\u07dc\x8d\u007fo\u0525\xbb\x93V[\xb2;\xef#h{!B\x01\x82\x02\xbf\xf9\xc3-\u0656[oO[`\xcc~\x8e\x8c*\xe0G\x13\xd9\xf8T\x94Node - us-central1-2\x9etaipei-us-central1-2@dexon.org\x8dus-central1-2\x91https://dexon.org\xf8\u0154\xa2!'\x95\x1c\xd44\xa7\u047aV\xe2\xa4\xf4.\xf3\xd7^]B\xf8\xae\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\x8cH3\x0eu\xac\xb0oM\xd8~\x87_\xb8\xecY\x0e\xd7\x03\x88\xe4B\x01\x8d<\x85[C{\x83\v\x16?k\x93\xfc\xca\xda`\u007f\x8e\xe2\x1d\xf1@\xe1M$\xaaB\xf6\xce\f\x1f1yQ\xea\x03\xec\x04\x99bX\xf8Q\x93Node - asia-east2-2\x9dtaipei-asia-east2-2@dexon.org\x8casia-east2-2\x91https://dexon.org\xf8\u0154\xa3\xd8\x1f\x16\xf7\x9c\x122\xa5BR\xbd\xb6\xa7^})%\x1a\xda\xf8\xae\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xf7I\xb9\xab\xa5.\x99r\xdc}\xf47\v\x0f~\x1c\u0557\xb8k\xee%\xde0\xdfs\x94P-\xc8\u0684\xf2\xf42E3\x1cy'\x90y\xe2\xa7\xef\xfa\xfd\xd2X\xbb{\x9d\t\xb5\x8e_\x17%\bq\xcb\xdaS\x02\xf8Q\x93Node - asia-east2-1\x9dtaipei-asia-east2-1@dexon.org\x8casia-east2-1\x91https://dexon.org\xf8\u0154\xa7<\x1c\xa6N\xf5\u007f|\xb9\x14.\xbd\xb4\x81:H\x97\xffh`\xf8\xae\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\u01e1\xaf\xc4%l\xa1\x87\xf2\xe8\x1eh\xad\xdf\xd7D\x8a\x1e\xa7\x15G\x8c/\x02\x8d/?\xd1\xf5\xf8\x9b\x82\x12\x06\x01y\xa0\r\x06\x11L \x0ecK\xecQ\x85\xbe\x01a\xf9,\xbcL\x9f\xfd\xdaI{\x83\x06L\xbf\xf8Q\x93Node - asia-east2-0\x9dtaipei-asia-east2-0@dexon.org\x8casia-east2-0\x91https://dexon.org\xf8\u0154\xae\x8b\x8b\xa2\xdfb}\xde_\xbb\x83\rn:\xe2D\x04\x19>\xaf\xf8\xae\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xd9z\xea\xb0K\x8d?\x8dH\xdc\xfa\x97\x9d\x04\x8b~\xb8b\x8b\xb8\xd6\xd7,3\xc5GGe\xe9\xc1\xbff\u02f3\xa2\xaa\xf0\x90\u049e\xcbEq\x92\xd7\x13\xc5~\xd9*\xf1\u0552|A,\xf3v^\x98!\xf4\x1d\x98\xf8Q\x93Node - asia-east1-0\x9dtaipei-asia-east1-0@dexon.org\x8casia-east1-0\x91https://dexon.org\xf8\u0214\xae\xd9Tzw\u074c\xabO\x1d8\xfaW\x19\"F\x17\xf8\x16E\xf8\xb1\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xe0\u80fb\xac\xe4\xd4\x10\xd3\a\xe3\xda0x\x17\xe6W\xe88\xb8:\x18)\xa2p\x19\x82\x00Ohr\xfd\x18\xb3\xc8\r%\xb4\u062b\x88P\xee#\u0570\x8b\x0e_\xa1\xb7\xac!\xa7\xaf\xe4\xd1y\x9e\x10\xad!VK\xf8T\x94Node - us-central1-0\x9etaipei-us-central1-0@dexon.org\x8dus-central1-0\x91https://dexon.org\xf8\u0514\xb8w\x97\x98-\x89\xe3z\xd8\u07a1t\xdex)\xaf\xc5\xfa\xb0\x8d\xf8\xbd\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xd8\xe9\xe2S\xe6\x04\x0f\x15\x13 \xea\x06\x82\xf0B\xc3 \xfcX\xf7\u05c5\xae\x83\x91g\xad\xb5\x1b\xbb\xd7\xe9C\xcbp0_\xf9h\x94~\xc2\xed\xa5[_\xd8i\xefU\x0e\xc0\x10\xf7'?W\xf7\x1c\x97\xadk=W\xf8`\x98Node - asia-northeast1-0\xa2taipei-asia-northeast1-0@dexon.org\x91asia-northeast1-0\x91https://dexon.org\xf8\u0154\xb8\u007fY\u0188\xeb\x0e\x97\x19\xbc\x91\xe8\xc0\xdb\x1c\x10\xd5P\x95j\xf8\xae\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\u007f\xf7i\xf9\xcc\t\x8ev4\xc2OIoO]/h\xd8\xeb\xd1*\xbaq\xdaY\xeb\x02(\x00o\xae8\x9b\u007fz\xff\xc9\xd1\x1e\x10\xae\x82\xff\x9c5\x9e\x92)\xb9\xefB=\xb3\u053b\u07f3\u057aj0[\xdcy\xf8Q\x93Node - asia-east1-1\x9dtaipei-asia-east1-1@dexon.org\x8casia-east1-1\x91https://dexon.org\xf8\xbf\x94\xb9\xa4wf\xca9\xe4O\xfd\U0007693a\u0165\xae1\xaf\xda\xe8\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04]\xa1\"\xd6\xfb9C\xf1V'\x99\x8d\xbb\x02C\xde|\x8a\x009\x93\xe4\x8a\xd71I2\xa6jo\xa2x\xa4\xf4\x94g\xd8\x15]\x9a%\xf1\x8fH&\xbe,\x99\x94\xa4\x95r\xc1\u03b4J\xff\x9b%j\x126,=\xf8K\x91Node - us-east4-2\x9btaipei-us-east4-2@dexon.org\x8aus-east4-2\x91https://dexon.org\ua53f\x8cH\xa6 \xba\xccF\x90\u007f\x9b\x89s-%\xe4z-|\xf7\u050b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x80\x80\x80\u0100\x80\x80\x80\xf8\xbf\x94\xd1\xef\xe1\xf9r\xcb-\u028b\u0241,\xc1\xb9\x91sv\x9b\x1c\x87\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04FU\xe6\a\xeb8\x95r\xe9\u6967\x88\xef\u061e\x94p\xc6\x15\x98o\x12\xc2\xf6l^\xe6F\xe2\x9eU'0\u0754\x88\xdct\xfb\xfd\u0741\xc7a\x8f\xe4\xae%\x05\xb8\xe6\x13hG\xaa\xe7yO\ub2bb\xfa\x01\xf8K\x91Node - us-east1-0\x9btaipei-us-east1-0@dexon.org\x8aus-east1-0\x91https://dexon.org\xf8\xbf\x94\xd3\x1d+#\x17@k\x16}\xde?NV\xd9x\x05q\x85GD\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\x8c0\x91\xf3\xbe\xa2&\x85@\x81\xbc\xbd\xcf\xd01\xc1\\]\xfd\xad\xd0ps\x92\xb2i\xce5#\x95\xbf[dQ\xa5\x8c\x1a8\xfc\xb5d\xfd?,\xb5\x0fLC\xf3{\xf8m#\xbd\xa5E^\xfdP\xdd_RW\xf1\xf8K\x91Node - us-west2-2\x9btaipei-us-west2-2@dexon.org\x8aus-west2-2\x91https://dexon.org\xf8\xec\x94\u0677r\b\xa5\xb7\xe4\xad5\xf9(\x80\xc53\xc2wV(\u07f7\xf8\u054b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04N\xf6\xd5t\xe1\xe8\xbe@[R\xc6\xe5\u0552\xb7cPvs\xd2+}\xa5\ueda2$\x89\x8b\xa5\xa8\xe9(\x1d\xecB\xf7]\xd4N\xb7g?D1h\x17\xc3\xf1\xaa\xc1N\xfd\x1a\xc2#\xc3\x1c\xa6h\x19\x84\x9c\t\xf8x\xa0Node - northamerica-northeast1-2\xaataipei-northamerica-northeast1-2@dexon.org\x99northamerica-northeast1-2\x91https://dexon.org\xf8\u02d4\u0713\x8c\xaf\x06\x8eE\xa0\x01\xa9Q.\x05D]V]\x19\xcc5\xf8\xb4\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xc2\xe5b\u03a2\xaa~\xf9\xa3t\v\xb9\xcf\\g\xb4\x1b\xb5\x11\x92'k\xd3\xf9n\xb8\u0278\xb3\xae,\x03\x1d^p\xac\u04df].k\x1f\xa2\xad!\xf2\x10\x0f\x85s\xa3nZfE\xf1\xf4j(\x8c\x87$*\xd8\xf8W\x95Node - europe-west4-0\x9ftaipei-europe-west4-0@dexon.org\x8eeurope-west4-0\x91https://dexon.org\xea\x94\xe0\xf8Y4\x03S\x85F\x93\xf57\xea3q\x06\xa3>\x9f\xea\xb0\u050b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x80\x80\x80\u0100\x80\x80\x80\xf8\u0214\xe1\xe5\x13\xfah\xeffr^\xda\xf8)\xf9AT<9\x13Kq\xf8\xb1\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\t\xe3\x0e\xfeq\x0e]\u06d3Z\xcc\n\xce^f\xa7\xb2\xf1r5!S\xe0p\xb4\x83=\xe4\x840\x8b\xab\xc0\x81\xe0\xe1\x84W$zx\xf9`cvB\x9a\x8f\x15\xc2T\x85\xe2\xbe\x00\u027cYk,\x8d\xa3*\x0e\xf8T\x94Node - us-central1-1\x9etaipei-us-central1-1@dexon.org\x8dus-central1-1\x91https://dexon.org\xf8\xec\x94\xef]\x14*\xf8\u04ecB$3\x17\xe7\xeb%\x8dY<\xabEu\xf8\u054b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04\xdf\xe2\x9c\xd7\x1e@Y\xbc\aD\v\x9d\xad\xedP\xa9\x1f\xf2\x14_\xa8.\u058e\x89\xabe\x1d\xcd7\x8dzL\xcdK\x92\xf5B\xfa?\bbw\xf2T\xbcW@Mx\xc3h\xb1\xa9\xe6\xabvZZv\x9e\x1e3?\xf8x\xa0Node - northamerica-northeast1-1\xaataipei-northamerica-northeast1-1@dexon.org\x99northamerica-northeast1-1\x91https://dexon.org\xf8\xbf\x94\xffh\xe5Yg\xc4\x15\xc3,\xb2\xa8K\xe6\xcd\xf4}U\xa6\xa7\u007f\xf8\xa8\x8b\x01\xa7\x847\x9d\x99\xdbB\x00\x00\x00\x8a\xd3\xc2\x1b\xce\xcc\xed\xa1\x00\x00\x00\x80\xb8A\x04w7\xd9\xc0\x9a=sJ\xb5i\xebg'\x85\x9f>#\x1b\xbf_\xef\x00=\x15\xdal\x95^\xcc#\a\x8e\x89C\xcc\x1e\xd9ji\xcc\x147\xaa\x93\xdb\x00_%>\\\xbe\x87\x8b3\xb0\v\xec\t\xaf\xdd\u03aa\x83\xc5\xf8K\x91Node - us-west2-0\x9btaipei-us-west2-0@dexon.org\x8aus-west2-0\x91https://dexon.org"
