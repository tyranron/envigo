// Copyright 2017 tyranron
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package envigo

import (
	"errors"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	Convey("Parse()", t, func() {
		Convey("Performs parsing with default parser", func() {
			setEnv("BOOL", "false")
			obj := &struct {
				V bool `env:"BOOL"`
			}{true}
			err := Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldBeFalse)
		})

		Convey("Returns default parser error if occurs", func() {
			obj1 := struct {
				V bool
			}{}
			some := true
			obj2 := &some

			So(Parse(obj1), ShouldEqual, ErrNotStructPtr)
			So(Parse(obj2), ShouldEqual, ErrNotStructPtr)
		})
	})
}

func TestParser_Parse(t *testing.T) {
	Convey("Parser.Parse()", t, func() {
		p := Parser{}

		Convey("If non-struct pointer is passed", func() {
			obj1 := struct {
				V bool
			}{}
			some := true
			obj2 := &some

			Convey("Returns error", func() {
				So(p.Parse(obj1), ShouldEqual, ErrNotStructPtr)
				So(p.Parse(obj2), ShouldEqual, ErrNotStructPtr)
			})
		})

		Convey("Parses supported types", func() {
			Convey("bool", func() {
				setEnv("BOOL", "false")
				obj := &struct {
					V bool `env:"BOOL"`
				}{true}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldBeFalse)
			})

			Convey("string", func() {
				setEnv("STRING", "foo")
				obj := &struct {
					V string `env:"STRING"`
				}{"bar"}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, "foo")
			})

			Convey("int", func() {
				setEnv("INT", "0")
				obj := &struct {
					V int `env:"INT"`
				}{-1}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, int(0))
			})

			Convey("int8", func() {
				setEnv("INT8", "-123")
				obj := &struct {
					V int8 `env:"INT8"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, int8(-123))
			})

			Convey("int16", func() {
				setEnv("INT16", "-32760")
				obj := &struct {
					V int16 `env:"INT16"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, int16(-32760))
			})

			Convey("int32", func() {
				setEnv("INT32", "-8388600")
				obj := &struct {
					V int32 `env:"INT32"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, int32(-8388600))
			})

			Convey("int64", func() {
				setEnv("INT64", "-2147483640")
				obj := &struct {
					V int64 `env:"INT64"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, int64(-2147483640))
			})

			Convey("uint", func() {
				setEnv("UINT", "0")
				obj := &struct {
					V uint `env:"UINT"`
				}{1}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, uint(0))
			})

			Convey("uint8", func() {
				setEnv("UINT8", "250")
				obj := &struct {
					V uint8 `env:"UINT8"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, uint8(250))
			})

			Convey("uint16", func() {
				setEnv("UINT16", "65530")
				obj := &struct {
					V uint16 `env:"UINT16"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, uint16(65530))
			})

			Convey("uint32", func() {
				setEnv("UINT32", "16777210")
				obj := &struct {
					V uint32 `env:"UINT32"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, uint32(16777210))
			})

			Convey("uint64", func() {
				setEnv("UINT64", "4294967290")
				obj := &struct {
					V uint64 `env:"UINT64"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, uint64(4294967290))
			})

			Convey("byte", func() {
				setEnv("BYTE", "255")
				obj := &struct {
					V byte `env:"BYTE"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, byte(255))
			})

			Convey("rune", func() {
				setEnv("RUNE", "8388600")
				obj := &struct {
					V rune `env:"RUNE"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, rune(8388600))
			})

			Convey("float32", func() {
				setEnv("FLOAT32", "3.40282346638528859811704183484516925440e+38")
				obj := &struct {
					V float32 `env:"FLOAT32"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual,
					float32(3.40282346638528859811704183484516925440e+38))
			})

			Convey("float64", func() {
				setEnv("FLOAT64", "1.797693134862315708145274237317043567981e+308")
				obj := &struct {
					V float64 `env:"FLOAT64"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual,
					float64(1.797693134862315708145274237317043567981e+308))
			})

			Convey("time.Duration", func() {
				setEnv("DURATION", "-1h2m3s4ms5us6ns")
				obj := &struct {
					V time.Duration `env:"DURATION"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual,
					-(time.Hour + 2*time.Minute + 3*time.Second +
						4*time.Millisecond + 5*time.Microsecond +
						6*time.Nanosecond))
			})

			Convey("time.Time", func() {
				t := time.Now()
				setEnv("TIME_RFC3339", t.Format(time.RFC3339))
				obj := &struct {
					V time.Time `env:"TIME_RFC3339"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V.Format(time.RFC3339), ShouldEqual, t.Format(time.RFC3339))
			})

			Convey("net.IP", func() {
				ipv4 := "32.1.219.8"
				ipv6 := "2001:db8:a0b:12f0::1"
				setEnv("IPv4", ipv4)
				setEnv("IPv6", ipv6)
				obj := &struct {
					V4 net.IP `env:"IPv4"`
					V6 net.IP `env:"IPv6"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V4.String(), ShouldEqual, ipv4)
				So(obj.V6.String(), ShouldEqual, ipv6)
			})

			Convey("array of", func() {
				Convey("bool", func() {
					setEnv("ARRAY_BOOL", "false,true,false")
					obj := &struct {
						V [3]bool `env:"ARRAY_BOOL"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [3]bool{false, true, false})
				})

				Convey("string", func() {
					arr := [3]string{"123.234.234:34", "34234:34234", "4234"}
					setEnv("ARRAY_STRING", strings.Join(arr[:], ","))
					obj := &struct {
						V [3]string `env:"ARRAY_STRING"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, arr)
				})

				Convey("int", func() {
					setEnv("ARRAY_INT", "0")
					obj := &struct {
						V [1]int `env:"ARRAY_INT"`
					}{[1]int{-1}}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]int{0})
				})

				Convey("int8", func() {
					setEnv("ARRAY_INT8", "-123")
					obj := &struct {
						V [1]int8 `env:"ARRAY_INT8"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]int8{-123})
				})

				Convey("int16", func() {
					setEnv("ARRAY_INT16", "-32760")
					obj := &struct {
						V [1]int16 `env:"ARRAY_INT16"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]int16{-32760})
				})

				Convey("int32", func() {
					setEnv("ARRAY_INT32", "-8388600")
					obj := &struct {
						V [1]int32 `env:"ARRAY_INT32"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]int32{-8388600})
				})

				Convey("int64", func() {
					setEnv("ARRAY_INT64", "-2147483640")
					obj := &struct {
						V [1]int64 `env:"ARRAY_INT64"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]int64{-2147483640})
				})

				Convey("uint", func() {
					setEnv("ARRAY_UINT", "0")
					obj := &struct {
						V [1]uint `env:"ARRAY_UINT"`
					}{[1]uint{1}}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]uint{0})
				})

				Convey("uint8", func() {
					setEnv("ARRAY_UINT8", "250")
					obj := &struct {
						V [1]uint8 `env:"ARRAY_UINT8"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]uint8{250})
				})

				Convey("uint16", func() {
					setEnv("ARRAY_UINT16", "65530")
					obj := &struct {
						V [1]uint16 `env:"ARRAY_UINT16"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]uint16{65530})
				})

				Convey("uint32", func() {
					setEnv("ARRAY_UINT32", "16777210")
					obj := &struct {
						V [1]uint32 `env:"ARRAY_UINT32"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]uint32{16777210})
				})

				Convey("uint64", func() {
					setEnv("ARRAY_UINT64", "4294967290")
					obj := &struct {
						V [1]uint64 `env:"ARRAY_UINT64"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]uint64{4294967290})
				})

				Convey("byte", func() {
					setEnv("ARRAY_BYTE", "255")
					obj := &struct {
						V [1]byte `env:"ARRAY_BYTE"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]byte{255})
				})

				Convey("rune", func() {
					setEnv("ARRAY_RUNE", "8388600")
					obj := &struct {
						V [1]rune `env:"ARRAY_RUNE"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]rune{8388600})
				})

				Convey("float32", func() {
					setEnv("ARRAY_FLOAT32",
						"3.40282346638528859811704183484516925440e+38")
					obj := &struct {
						V [1]float32 `env:"ARRAY_FLOAT32"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]float32{
						3.40282346638528859811704183484516925440e+38,
					})
				})

				Convey("float64", func() {
					setEnv("ARRAY_FLOAT64",
						"1.797693134862315708145274237317043567981e+308")
					obj := &struct {
						V [1]float64 `env:"ARRAY_FLOAT64"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [1]float64{
						1.797693134862315708145274237317043567981e+308,
					})
				})

				Convey("time.Duration", func() {
					setEnv("ARRAY_DURATION", "-1h2m3s4ms5us6ns,3h")
					obj := &struct {
						V [2]time.Duration `env:"ARRAY_DURATION"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, [2]time.Duration{
						-(time.Hour + 2*time.Minute + 3*time.Second +
							4*time.Millisecond + 5*time.Microsecond +
							6*time.Nanosecond),
						3 * time.Hour,
					})
				})

				Convey("net.IP", func() {
					setEnv("ARRAY_IP", "32.1.219.8,2001:db8:a0b:12f0::1")
					obj := &struct {
						V [2]net.IP `env:"ARRAY_IP"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V[0].String(), ShouldEqual, "32.1.219.8")
					So(obj.V[1].String(), ShouldEqual, "2001:db8:a0b:12f0::1")
				})
			})

			Convey("slice of", func() {
				Convey("bool", func() {
					setEnv("SLICE_BOOL", "false,true,false")
					obj := &struct {
						V []bool `env:"SLICE_BOOL"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []bool{false, true, false})
				})

				Convey("string", func() {
					slice := []string{"123.234.234:34", "34234:34234", "4234"}
					setEnv("SLICE_STRING", strings.Join(slice, ","))
					obj := &struct {
						V []string `env:"SLICE_STRING"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, slice)
				})

				Convey("int", func() {
					setEnv("SLICE_INT", "0")
					obj := &struct {
						V []int `env:"SLICE_INT"`
					}{[]int{-1}}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []int{0})
				})

				Convey("int8", func() {
					setEnv("SLICE_INT8", "-123")
					obj := &struct {
						V []int8 `env:"SLICE_INT8"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []int8{-123})
				})

				Convey("int16", func() {
					setEnv("SLICE_INT16", "-32760")
					obj := &struct {
						V []int16 `env:"SLICE_INT16"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []int16{-32760})
				})

				Convey("int32", func() {
					setEnv("SLICE_INT32", "-8388600")
					obj := &struct {
						V []int32 `env:"SLICE_INT32"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []int32{-8388600})
				})

				Convey("int64", func() {
					setEnv("SLICE_INT64", "-2147483640")
					obj := &struct {
						V []int64 `env:"SLICE_INT64"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []int64{-2147483640})
				})

				Convey("uint", func() {
					setEnv("SLICE_UINT", "0")
					obj := &struct {
						V []uint `env:"SLICE_UINT"`
					}{[]uint{1}}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []uint{0})
				})

				Convey("uint8", func() {
					setEnv("SLICE_UINT8", "250")
					obj := &struct {
						V []uint8 `env:"SLICE_UINT8"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []uint8{250})
				})

				Convey("uint16", func() {
					setEnv("SLICE_UINT16", "65530")
					obj := &struct {
						V []uint16 `env:"SLICE_UINT16"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []uint16{65530})
				})

				Convey("uint32", func() {
					setEnv("SLICE_UINT32", "16777210")
					obj := &struct {
						V []uint32 `env:"SLICE_UINT32"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []uint32{16777210})
				})

				Convey("uint64", func() {
					setEnv("SLICE_UINT64", "4294967290")
					obj := &struct {
						V []uint64 `env:"SLICE_UINT64"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []uint64{4294967290})
				})

				Convey("byte", func() {
					setEnv("SLICE_BYTE", "255")
					obj := &struct {
						V []byte `env:"SLICE_BYTE"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []byte{255})
				})

				Convey("rune", func() {
					setEnv("SLICE_RUNE", "8388600")
					obj := &struct {
						V []rune `env:"SLICE_RUNE"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []rune{8388600})
				})

				Convey("float32", func() {
					setEnv("SLICE_FLOAT32",
						"3.40282346638528859811704183484516925440e+38")
					obj := &struct {
						V []float32 `env:"SLICE_FLOAT32"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []float32{
						3.40282346638528859811704183484516925440e+38,
					})
				})

				Convey("float64", func() {
					setEnv("SLICE_FLOAT64",
						"1.797693134862315708145274237317043567981e+308")
					obj := &struct {
						V []float64 `env:"SLICE_FLOAT64"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []float64{
						1.797693134862315708145274237317043567981e+308,
					})
				})

				Convey("time.Duration", func() {
					setEnv("SLICE_DURATION", "-1h2m3s4ms5us6ns,3h")
					obj := &struct {
						V []time.Duration `env:"SLICE_DURATION"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldResemble, []time.Duration{
						-(time.Hour + 2*time.Minute + 3*time.Second +
							4*time.Millisecond + 5*time.Microsecond +
							6*time.Nanosecond),
						3 * time.Hour,
					})
				})

				Convey("net.IP", func() {
					setEnv("SLICE_IP", "32.1.219.8,2001:db8:a0b:12f0::1")
					obj := &struct {
						V []net.IP `env:"SLICE_IP"`
					}{}
					err := p.Parse(obj)

					So(err, ShouldBeNil)
					So(obj.V, ShouldHaveLength, 2)
					So(obj.V[0].String(), ShouldEqual, "32.1.219.8")
					So(obj.V[1].String(), ShouldEqual, "2001:db8:a0b:12f0::1")
				})
			})
		})

		Convey("On unsupported type", func() {
			setEnv("UNSUPPORTED_TYPE", "2")
			obj := &struct {
				V uintptr `env:"UNSUPPORTED_TYPE"`
			}{5}
			err := p.Parse(obj)

			Convey("Returns error", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldHaveSameTypeAs, UnparsableTypeError{})
			})

			Convey("Does not mutate value", func() {
				So(obj.V, ShouldEqual, uintptr(5))
			})
		})

		Convey("On incorrectly declared tag", func() {
			setEnv("UINT8", "3")
			obj := &struct {
				V uint8 `env:""`
			}{5}
			err := p.Parse(obj)

			Convey("Returns error", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldHaveSameTypeAs, EmptyVarNameError{})
			})

			Convey("Does not mutate value", func() {
				So(obj.V, ShouldEqual, 5)
			})
		})

		Convey("If env var is not set", func() {
			Convey("Does not mutate value", func() {
				unsetEnv("UINT8")
				obj := &struct {
					V uint8 `env:"UINT8"`
				}{5}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, 5)
			})
		})

		Convey("If env var is empty", func() {
			Convey("Parses empty value", func() {
				setEnv("STRING", "")
				obj := &struct {
					V string `env:"STRING"`
				}{"some"}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, "")
			})
		})

		Convey("Parses nested structs", func() {
			setEnv("NESTED_BOOL", "true")
			setEnv("NESTED_INT", "-10")
			setEnv("NESTED_UINT", "15")
			obj := &struct {
				V bool `env:"NESTED_BOOL"`
				N struct {
					V int `env:"NESTED_INT"`
					N struct {
						N struct {
							V uint `env:"NESTED_UINT"`
						}
					}
				}
				h bool // nolint: unused, megacheck
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldEqual, true)
			So(obj.N.V, ShouldEqual, int(-10))
			So(obj.N.N.N.V, ShouldEqual, uint(15))
		})

		Convey("Parses embedded structs", func() {
			setEnv("EMBEDDED_BOOL", "true")
			setEnv("EMBEDDED_INT", "-2")

			Convey("Raw struct", func() {
				obj := &struct {
					EmbeddedStruct
					V bool `env:"EMBEDDED_BOOL"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, true)
				So(obj.EmbeddedStruct.V, ShouldEqual, true)
				So(obj.EmbeddedStruct.V2, ShouldEqual, -2)
			})

			Convey("Struct behind pointer", func() {
				obj := &struct {
					*EmbeddedStruct
					V bool `env:"EMBEDDED_BOOL"`
				}{EmbeddedStruct: &EmbeddedStruct{}}
				err := p.Parse(obj)

				So(err, ShouldBeNil)
				So(obj.V, ShouldEqual, true)
				So(obj.EmbeddedStruct.V, ShouldEqual, true)
				So(obj.EmbeddedStruct.V2, ShouldEqual, -2)
			})
		})

		Convey("Parses values for types behind pointers", func() {
			setEnv("DEREF_BOOL", "true")
			setEnv("DEREF_INT", "-10")
			i := 5
			ptr1 := &i
			ptr2 := &ptr1
			b := true
			obj := &struct {
				V *bool `env:"DEREF_BOOL"`
				N *struct {
					V ***int `env:"DEREF_INT"`
				}
			}{V: &b, N: &struct {
				V ***int `env:"DEREF_INT"`
			}{&ptr2}}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(*(obj.V), ShouldEqual, true)
			So(***(obj.N).V, ShouldEqual, int(-10))
		})

		Convey("Omitts nil pointers", func() {
			setEnv("PTR_BOOL", "true")
			obj := &struct {
				V *bool `env:"PTR_BOOL"`
			}{}
			err := p.Parse(obj)

			So(err, ShouldBeNil)
			So(obj.V, ShouldBeNil)
		})

		Convey("Uses custom parser if type has one", func() {
			Convey("Performs custom parse correctly", func() {
				setEnv("CUSTOM_UINT8", "10")
				obj1 := &struct {
					V customUint8 `env:"CUSTOM_UINT8"`
				}{}
				v2 := customUint8(2)
				obj2 := &struct {
					V *customUint8 `env:"CUSTOM_UINT8"`
				}{&v2}
				v3 := customUint8(3)
				pv3 := &v3
				obj3 := &struct {
					V **customUint8 `env:"CUSTOM_UINT8"`
				}{&pv3}
				err1 := p.Parse(obj1)
				err2 := p.Parse(obj2)
				err3 := p.Parse(obj3)

				So(err1, ShouldBeNil)
				So(obj1.V, ShouldEqual, 7)
				So(err2, ShouldBeNil)
				So(*(obj2.V), ShouldEqual, 7)
				So(err3, ShouldBeNil)
				So(**(obj3.V), ShouldEqual, 7)
			})
		})

		Convey("On tagged struct without custom parser", func() {
			Convey("Returns error of unsupported type", func() {
				setEnv("EMBEDDED_STRUCT", "{3}")
				setEnv("EMBEDDED_BOOL", "true")
				setEnv("EMBEDDED_INT", "-2")
				obj := &struct {
					V EmbeddedStruct `env:"EMBEDDED_STRUCT"`
				}{}
				err := p.Parse(obj)

				So(err, ShouldNotBeNil)
				So(err, ShouldHaveSameTypeAs, UnparsableTypeError{})
			})
		})

		Convey("If value cannot be parsed", func() {
			Convey("Returns parsing error", func() {
				Convey("supported types", func() {
					setEnv("FAIL_BOOL", "hi")
					setEnv("FAIL_INT", "true")
					setEnv("FAIL_UINT", "false")
					setEnv("FAIL_FLOAT", "-----")
					setEnv("FAIL_DURATION", "???")
					obj1 := &struct {
						V bool `env:"FAIL_BOOL"`
					}{}
					obj2 := &struct {
						V int `env:"FAIL_INT"`
					}{}
					obj3 := &struct {
						V uint16 `env:"FAIL_UINT"`
					}{}
					obj4 := &struct {
						V float64 `env:"FAIL_FLOAT"`
					}{}
					obj5 := &struct {
						V time.Duration `env:"FAIL_DURATION"`
					}{}
					err1 := p.Parse(obj1)
					err2 := p.Parse(obj2)
					err3 := p.Parse(obj3)
					err4 := p.Parse(obj4)
					err5 := p.Parse(obj5)

					So(err1, ShouldNotBeNil)
					So(err1, ShouldHaveSameTypeAs, ParseError{})
					So(err1.Error(), ShouldContainSubstring, "'FAIL_BOOL'")
					So(err2, ShouldNotBeNil)
					So(err2, ShouldHaveSameTypeAs, ParseError{})
					So(err2.Error(), ShouldContainSubstring, "'FAIL_INT'")
					So(err3, ShouldNotBeNil)
					So(err3, ShouldHaveSameTypeAs, ParseError{})
					So(err3.Error(), ShouldContainSubstring, "'FAIL_UINT'")
					So(err4, ShouldNotBeNil)
					So(err4, ShouldHaveSameTypeAs, ParseError{})
					So(err4.Error(), ShouldContainSubstring, "'FAIL_FLOAT'")
					So(err5, ShouldNotBeNil)
					So(err5, ShouldHaveSameTypeAs, ParseError{})
					So(err5.Error(), ShouldContainSubstring, "'FAIL_DURATION'")
				})

				Convey("custom parser type", func() {
					setEnv("FAIL_CUSTOM", "10")
					v := customFailure(4)
					pV := &v
					obj1 := &struct {
						V *customFailure `env:"FAIL_CUSTOM"`
					}{pV}
					obj2 := &struct {
						V **customFailure `env:"FAIL_CUSTOM"`
					}{&pV}
					err1 := p.Parse(obj1)
					err2 := p.Parse(obj2)

					So(err1, ShouldNotBeNil)
					So(err1, ShouldHaveSameTypeAs, ParseError{})
					So(err2, ShouldNotBeNil)
					So(err2, ShouldHaveSameTypeAs, ParseError{})
				})

				Convey("nested structs", func() {
					setEnv("FAIL_NESTED", "?-?")
					obj1 := &struct {
						N struct {
							V int `env:"FAIL_NESTED"`
						}
					}{}
					obj2 := &struct {
						N *struct {
							V int `env:"FAIL_NESTED"`
						}
					}{&struct {
						V int `env:"FAIL_NESTED"`
					}{}}
					err1 := p.Parse(obj1)
					err2 := p.Parse(obj2)

					So(err1, ShouldNotBeNil)
					So(err1, ShouldHaveSameTypeAs, ParseError{})
					So(err2, ShouldNotBeNil)
					So(err2, ShouldHaveSameTypeAs, ParseError{})
				})
			})
		})
	})
}

type customUint8 uint8

func (v *customUint8) UnmarshalText(_ []byte) error {
	*v = 7
	return nil
}

type customFailure uint8

func (_ *customFailure) UnmarshalText(_ []byte) error {
	return errors.New("some error")
}

type EmbeddedStruct struct {
	V  bool `env:"EMBEDDED_BOOL"`
	V2 int  `env:"EMBEDDED_INT"`
}

// setEnv is a simple helper function for setting env vars in one line.
func setEnv(name, val string) {
	if err := os.Setenv(name, val); err != nil {
		panic(err)
	}
}

// unsetEnv is a simple helper function for removing env vars in one line.
func unsetEnv(name string) {
	if err := os.Unsetenv(name); err != nil {
		panic(err)
	}
}
