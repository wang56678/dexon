package runtime

var testData = &tmplData{
	BinOpCollections: []*tmplTestCollection{
		{
			TestName: "OpAdd", OpFunc: "opAdd",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "ADD",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data: []string{
								"{V: 1, V: 2}", "{V: -1, V: -2}", "{V: 10, V: 10}", "{V: -10, V: 10}",
							},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data: []string{
								"{V: 1, V: 2}", "{V: 1, V: 2}", "{V: -10, V: 10}", "{V: -10, V: 3}",
							},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
						Data: []string{
							"{V: 2, V: 4}", "{V: 0, V: 0}", "{V: 0, V: 20}", "{V: -20, V: 13}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "ADD",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data: []string{
								"{V: 10, V: 10, V: 10}", "{V: -10, V: -10, V: -10}", "{V: -1, V: 2, V: 0}",
							},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data: []string{
								"{V: 1, V: -2, V: -10}",
							},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
						Data: []string{
							"{V: 11, V: 8, V: 0}", "{V: -9, V: -12, V: -20}", "{V: 0, V: 0, V: -10}",
						},
					},
				},
				{
					Name:  "Immediate 2",
					Error: "nil", OpCode: "ADD",
					Inputs: []*tmplOp{
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 1, V: -2, V: -10}",
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 10, V: 10, V: 10}",
								"{V: -10, V: -10, V: -10}",
								"{V: -1, V: 2, V: 0}",
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 0},
							{Major: "Int", Minor: 0},
							{Major: "Int", Minor: 0},
						},
						Data: []string{
							"{V: 11, V: 8, V: 0}",
							"{V: -9, V: -12, V: -20}",
							"{V: 0, V: 0, V: -10}",
						},
					},
				},
				{
					Name:  "Overflow - Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "ADD",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 127}", "{V: 1}", "{V: 1}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Overflow None Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "ADD",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 126}", "{V: 126}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}", "{V: 2}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Underflow - Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "ADD",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -128}", "{V: -1}", "{V: -1}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -1}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Underflow None Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "ADD",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -127}", "{V: -127}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -1}", "{V: -2}"},
						},
					},
					Output: &tmplOp{},
				},
			},
		},
		// -- end of ADD
		{
			TestName: "OpSub", OpFunc: "opSub",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "SUB",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 2}", "{V: -1, V: -2}", "{V: 10, V: 10}", "{V: -10, V: 10}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 2}", "{V: 1, V: 2}", "{V: -10, V: 10}", "{V: 10, V: 3}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
						Data:  []string{"{V: 0, V: 0}", "{V: -2, V: -4}", "{V: 20, V: 0}", "{V: -20, V: 7}"},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "SUB",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data: []string{
								"{V: 10, V: 10, V: 10}", "{V: -10, V: -10, V: -10}", "{V: -1, V: 2, V: 0}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: -2, V: -10}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
						Data:  []string{"{V: 9, V: 12, V: 20}", "{V: -11, V: -8, V: 0}", "{V: -2, V: 4, V: 10}"},
					},
				},
				{
					Name:  "Immediate 2",
					Error: "nil", OpCode: "SUB",
					Inputs: []*tmplOp{
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 1, V: -2, V: -10}",
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 10, V: 10, V: 10}",
								"{V: -10, V: -10, V: -10}",
								"{V: -1, V: 2, V: 0}",
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 0},
							{Major: "Int", Minor: 0},
							{Major: "Int", Minor: 0},
						},
						Data: []string{
							"{V: -9, V: -12, V: -20}",
							"{V: 11, V: 8, V: 0}",
							"{V: 2, V: -4, V: -10}",
						},
					},
				},
				{
					Name:  "Overflow - Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "SUB",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 127}", "{V: 1}", "{V: 1}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -1}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Overflow None Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "SUB",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 126}", "{V: 126}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -1}", "{V: -2}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Underflow - Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "SUB",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -128}", "{V: -1}", "{V: -1}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Underflow None Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "SUB",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -127}", "{V: -127}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}", "{V: 2}"},
						},
					},
					Output: &tmplOp{},
				},
			},
		},
		// -- end of SUB
		{
			TestName: "OpMul", OpFunc: "opMul",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "MUL",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 1}", "{V: 2, V: -1}", "{V: -2, V: 10}", "{V: 10, V: -10}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 1}", "{V: 2, V: 1}", "{V: 2, V: -10}", "{V: 10, V: -10}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
						Data:  []string{"{V: 0, V: 1}", "{V: 4, V: -1}", "{V: -4, V: -100}", "{V: 100, V: 100}"},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "MUL",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 10, V: 10, V: 10}", "{V: -10, V: -10, V: -10}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: -2, V: 0}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
						Data:  []string{"{V: 10, V: -20, V: 0}", "{V: -10, V: 20, V: 0}"},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "MUL",
					Inputs: []*tmplOp{
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 1, V: -2, V: 0}",
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 10, V: 10, V: 10}",
								"{V: -10, V: -10, V: -10}",
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 0},
							{Major: "Int", Minor: 0},
							{Major: "Int", Minor: 0},
						},
						Data: []string{
							"{V: 10, V: -20, V: 0}",
							"{V: -10, V: 20, V: 0}",
						},
					},
				},
				{
					Name:  "Overflow - Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "MUL",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 127}", "{V: 1}", "{V: 1}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 2}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Overflow None Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "MUL",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 126}", "{V: 126}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}", "{V: 2}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Underflow - Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "MUL",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -128}", "{V: -1}", "{V: -1}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 2}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Underflow None Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "MUL",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -127}", "{V: -127}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}", "{V: 2}"},
						},
					},
					Output: &tmplOp{},
				},
			},
		},
		// -- end of MUL
		{
			TestName: "OpDiv", OpFunc: "opDiv",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "DIV",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 1}", "{V: 2, V: -1}", "{V: -2, V: 10}", "{V: 10, V: -10}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 1}", "{V: 2, V: 1}", "{V: 2, V: -10}", "{V: 10, V: -10}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
						Data:  []string{"{V: 0, V: 1}", "{V: 1, V: -1}", "{V: -1, V: -1}", "{V: 1, V: 1}"},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "DIV",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 0}", "{V: 10, V: 10}", "{V: -10, V: -10}", "{V: 13, V: 13}", "{V: -13, V: -13}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 2, V: -2}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
						Data:  []string{"{V: 0, V: 0}", "{V: 5, V: -5}", "{V: -5, V: 5}", "{V: 6, V: -6}", "{V: -6, V: 6}"},
					},
				},
				{
					Name:  "Immediate 2",
					Error: "nil", OpCode: "DIV",
					Inputs: []*tmplOp{
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 50, V: -50}",
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 10, V: 10}",
								"{V: -10, V: -10}",
								"{V: 9, V: 9}",
								"{V: -9, V: -9}",
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 0},
							{Major: "Int", Minor: 0},
						},
						Data: []string{
							"{V: 5, V: -5}",
							"{V: -5, V: 5}",
							"{V: 5, V: -5}",
							"{V: -5, V: 5}",
						},
					},
				},
				{
					Name:  "DivideByZero Immediate",
					Error: "errors.ErrorCodeDividedByZero", OpCode: "DIV",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 10}", "{V: 10}", "{V: 10}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "DivideByZero None Immediate",
					Error: "errors.ErrorCodeDividedByZero", OpCode: "DIV",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 10}", "{V: 10}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}", "{V: 0}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Overflow - Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "DIV",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}", "{V: -128}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -1}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Overflow None Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "DIV",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -128}", "{V: -128}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: -1}", "{V: -2}"},
						},
					},
					Output: &tmplOp{},
				},
			},
		},
		// -- end of DIV
		{
			TestName: "OpMod", OpFunc: "opMod",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "MOD",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 1}", "{V: 0, V: -1}", "{V: 2, V: -2}", "{V: 2, V: -2}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 1}", "{V: -1, V: -1}", "{V: 3, V: 3}", "{V: -3, V: -3}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
						Data:  []string{"{V: 0, V: 0}", "{V: 0, V: 0}", "{V: 2, V: -2}", "{V: 2, V: -2}"},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "MOD",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 0}", "{V: 10, V: 10}", "{V: -10, V: -10}", "{V: 13, V: 13}", "{V: -13, V: -13}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 3, V: -3}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
						Data:  []string{"{V: 0, V: 0}", "{V: 1, V: 1}", "{V: -1, V: -1}", "{V: 1, V: 1}", "{V: -1, V: -1}"},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "MOD",
					Inputs: []*tmplOp{
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 31, V: -31}",
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 10, V: 10}",
								"{V: -10, V: -10}",
								"{V: 13, V: 13}",
								"{V: -13, V: -13}",
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 0},
							{Major: "Int", Minor: 0},
						},
						Data: []string{
							"{V: 1, V: -1}",
							"{V: 1, V: -1}",
							"{V: 5, V: -5}",
							"{V: 5, V: -5}",
						},
					},
				},
				{
					Name:  "ModideByZero Immediate",
					Error: "errors.ErrorCodeDividedByZero", OpCode: "MOD",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 10}", "{V: 10}", "{V: 10}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "ModideByZero None Immediate",
					Error: "errors.ErrorCodeDividedByZero", OpCode: "MOD",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 10}", "{V: 10}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}", "{V: 0}"},
						},
					},
					Output: &tmplOp{},
				},
			},
		},
		// -- end of MOD
		{
			TestName: "OpConcat", OpFunc: "opConcat",
			Cases: []*tmplTestCase{
				{
					Name:  "Concat bytes",
					Error: "nil", OpCode: "CONCAT",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
							},
							Data: []string{
								`{B: "abc-1", B: "xyz-1"}`,
								`{B: "abc-2", B: "xyz-2"}`,
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
							},
							Data: []string{
								`{B: "ABC-1", B: "XYZ-1"}`,
								`{B: "ABC-2", B: "XYZ-2"}`,
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "DynamicBytes", Minor: 0},
							{Major: "DynamicBytes", Minor: 0},
						},
						Data: []string{
							`{B: "abc-1ABC-1", B: "xyz-1XYZ-1"}`,
							`{B: "abc-2ABC-2", B: "xyz-2XYZ-2"}`,
						},
					},
				},
				{
					Name:  "Invalid concat",
					Error: "errors.ErrorCodeInvalidDataType", OpCode: "CONCAT",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{`{B: "abc-1", T}`},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{`{B: "ABC-1", F}`},
						},
					},
					Output: &tmplOp{},
				},
			},
		},
		// -- end of CONCAT
		{
			TestName: "OpLt", OpFunc: "opLt",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "LT",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 0, V: -1}", "{V: 1, V: 0, V: -1}", "{V: 1, V: 0, V: -1}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 1, V: 1}", "{V: 0, V: 0, V: 0}", "{V: -1, V: -1, V: -1}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{F, T, T}", "{F, F, T}", "{F, F, F}"},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "LT",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 1, V: -1}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 0, V: 0}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{F, F, T}"},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "LT",
					Inputs: []*tmplOp{
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 0, V: 0, V: 0}",
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 0, V: 1, V: -1}",
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Bool", Minor: 0},
							{Major: "Bool", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							"{F, T, F}",
						},
					},
				},
			},
		},
		// -- end of LT
		{
			TestName: "OpGt", OpFunc: "opGt",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "GT",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 0, V: -1}", "{V: 1, V: 0, V: -1}", "{V: 1, V: 0, V: -1}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 1, V: 1}", "{V: 0, V: 0, V: 0}", "{V: -1, V: -1, V: -1}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{F, F, F}", "{T, F, F}", "{T, T, F}"},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "GT",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 1, V: -1}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 0, V: 0}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{F, T, F}"},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "GT",
					Inputs: []*tmplOp{
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 0, V: 0, V: 0}",
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								"{V: 0, V: 1, V: -1}",
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Bool", Minor: 0},
							{Major: "Bool", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							"{F, F, T}",
						},
					},
				},
			},
		},
		// -- end of GT
		{
			TestName: "OpEq", OpFunc: "opEq",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "EQ",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 0, V: -1}", "{V: 1, V: 0, V: -1}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 0, V: -1}", "{V: 1, V: 1, V: 1}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{T, T, T}", "{T, F, F}"},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "EQ",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 0, V: 0}", "{V: 0, V: 1, V: -1}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0, V: 0, V: 0}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{T, T, T}", "{T, F, F}"},
					},
				},
			},
		},
		// -- end of EQ
		{
			TestName: "OpAnd", OpFunc: "opAnd",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "AND",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, F}", "{F, T}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, T}", "{F, F}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{T, F}", "{F, F}"},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "AND",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, F}", "{F, T}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, T}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{T, F}", "{F, T}"},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "AND",
					Inputs: []*tmplOp{
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Bool", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								"{T, T}",
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Bool", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								"{T, F}",
								"{F, T}",
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Bool", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							"{T, F}",
							"{F, T}",
						},
					},
				},
				{
					Name:  "Invalid Data Type",
					Error: "errors.ErrorCodeInvalidDataType", OpCode: "AND",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}"},
						},
					},
					Output: &tmplOp{},
				},
			},
		},
		// -- end of AND
		{
			TestName: "OpOr", OpFunc: "opOr",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "OR",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, F}", "{F, T}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, T}", "{F, F}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{T, T}", "{F, T}"},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "OR",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, F}", "{F, T}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, T}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{T, T}", "{T, T}"},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "OR",
					Inputs: []*tmplOp{
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Bool", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								"{T, T}",
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Bool", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								"{T, F}",
								"{F, T}",
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Bool", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							"{T, T}",
							"{T, T}",
						},
					},
				},
				{
					Name:  "Invalid Data Type",
					Error: "errors.ErrorCodeInvalidDataType", OpCode: "OR",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}"},
						},
					},
					Output: &tmplOp{},
				},
			},
		},
		// -- end of OR
		{
			TestName: "OpNot", OpFunc: "opNot",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "NOT",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, F}", "{F, T}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{F, T}", "{T, F}"},
					},
				},
				{
					Name:  "Errors Invalid Data Type",
					Error: "errors.ErrorCodeInvalidDataType", OpCode: "NOT",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}"},
						},
					},
					Output: &tmplOp{},
				},
			},
		},
		// -- end of NOT
		{
			TestName: "OpUnion", OpFunc: "opUnion",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "UNION",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, F}", "{F, T}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, T}", "{F, F}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{F, F}", "{F, T}", "{T, F}", "{T, T}"},
					},
				},
			},
		},
		// -- end of UNION
		{
			TestName: "OpIntxn", OpFunc: "opIntxn",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "INTXN",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, F}", "{F, T}", "{T, T}", "{F, F}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, T}", "{F, F}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{F, F}", "{T, T}"},
					},
				},
			},
		},
		// -- end of INTXN
		{
			TestName: "OpLike", OpFunc: "opLike",
			Cases: []*tmplTestCase{
				{
					Name:  `Like %\\%b% escape \\`, // \\ is raw string escape for \
					Error: "nil", OpCode: "LIKE",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
							},
							Data: []string{
								`{B: "a%bcdefg", B: "gfedcba"}`,
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
							},
							Data: []string{`{B: "%\\%b%"}`},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
							},
							Data: []string{`{B: "\\"}`},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{T, F}"},
					},
				},
				{
					Name:  `Like t1 escape t2`,
					Error: "nil", OpCode: "LIKE",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
							},
							Data: []string{
								`{B: "a%bcdefg"}`,
								`{B: "gfedcba"}`,
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
							},
							Data: []string{
								`{B: "%\\%b%"}`,
								`{B: "_fed%"}`,
							},
						},
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
							},
							Data: []string{
								`{B: "\\"}`,
								`{B: ""}`,
							},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							"{T}",
							"{T}",
						},
					},
				},
				{
					Name:  `Like with valid and invalid UTF8`,
					Error: "nil", OpCode: "LIKE",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
							},
							Data: []string{
								`{B: {226, 40, 161, 228, 189, 160, 229, 165, 189}, B: "gfedcba"}`,
								// "\xe2(\xa1你好"
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
							},
							Data: []string{
								`{B: {37, 228, 189, 160, 37}}`,
								// "%你%"
							},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{T, F}"},
					},
				},
			},
		},
		// -- end of LIKE
		{
			TestName: "OpZip", OpFunc: "opZip",
			Cases: []*tmplTestCase{
				{
					Name:  "Zip two array",
					Error: "nil", OpCode: "ZIP",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "DynamicBytes", Minor: 0}, {Major: "DynamicBytes", Minor: 0}},
							Data:  []string{`{B: "abcdefg-1", B: "gfedcba-1"}`, `{B: "abcdefg-2", B: "gfedcba-2"}`},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{V: 1, T}", "{V: 2, F}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "DynamicBytes", Minor: 0},
							{Major: "DynamicBytes", Minor: 0},
							{Major: "Int", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							`{B: "abcdefg-1", B: "gfedcba-1", V: 1, T}`,
							`{B: "abcdefg-2", B: "gfedcba-2", V: 2, F}`,
						},
					},
				},
				{
					Name:  "Zip immediate",
					Error: "nil", OpCode: "ZIP",
					Inputs: []*tmplOp{
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "DynamicBytes", Minor: 0}, {Major: "DynamicBytes", Minor: 0}},
							Data:  []string{`{B: "abcdefg-1", B: "gfedcba-1"}`},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{V: 1, T}", "{V: 2, F}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "DynamicBytes", Minor: 0},
							{Major: "DynamicBytes", Minor: 0},
							{Major: "Int", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							`{B: "abcdefg-1", B: "gfedcba-1", V: 1, T}`,
							`{B: "abcdefg-1", B: "gfedcba-1", V: 2, F}`,
						},
					},
				},
			},
		},
		// -- end of ZIP
		{
			TestName: "OpField", OpFunc: "opField",
			Cases: []*tmplTestCase{
				{
					Name:  "Retrieve 2nd,3rd column",
					Error: "nil", OpCode: "FIELD",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								`{B: "abcdefg-1", B: "gfedcba-1", V: 1, T}`,
								`{B: "abcdefg-2", B: "gfedcba-2", V: 2, F}`,
							},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 2}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "DynamicBytes", Minor: 0}, {Major: "Int", Minor: 0}},
						Data:  []string{`{B: "gfedcba-1", V: 1}`, `{B: "gfedcba-2", V: 2}`},
					},
				},
			},
		},
		// -- end of FIELD
		{
			TestName: "OpPrune", OpFunc: "opPrune",
			Cases: []*tmplTestCase{
				{
					Name:  "Prune 2nd,4th,5th column",
					Error: "nil", OpCode: "PRUNE",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Bool", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								`{B: "abcdefg-1", B: "gfedcba-1", V: 1, F, T}`,
								`{B: "abcdefg-2", B: "gfedcba-2", V: 2, T, F}`,
							},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}, {Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1, V: 3, V: 4}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "DynamicBytes", Minor: 0}, {Major: "Int", Minor: 0}},
						Data:  []string{`{B: "abcdefg-1", V: 1}`, `{B: "abcdefg-2", V: 2}`},
					},
				},
			},
		},
		// -- end of PRUNE
		{
			TestName: "OpCut", OpFunc: "opCut",
			Cases: []*tmplTestCase{
				{
					Name:  "Cut 2nd to 4th columns",
					Error: "nil", OpCode: "CUT",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Bool", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								`{B: "abcdefg-1", B: "gfedcba-1", V: 1, F, T}`,
								`{B: "abcdefg-2", B: "gfedcba-2", V: 2, T, F}`,
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{"{V: 1, V: 3}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "DynamicBytes", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							`{B: "abcdefg-1", T}`,
							`{B: "abcdefg-2", F}`,
						},
					},
				},
				{
					Name:  "Cut 1st column",
					Error: "nil", OpCode: "CUT",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Bool", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								`{B: "abcdefg-1", B: "gfedcba-1", V: 1, F, T}`,
								`{B: "abcdefg-2", B: "gfedcba-2", V: 2, T, F}`,
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{"{V: 0, V: 0}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "DynamicBytes", Minor: 0},
							{Major: "Int", Minor: 0},
							{Major: "Bool", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							`{B: "gfedcba-1", V: 1, F, T}`,
							`{B: "gfedcba-2", V: 2, T, F}`,
						},
					},
				},
				{
					Name:  "Cut since 2nd column",
					Error: "nil", OpCode: "CUT",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								`{B: "abcdefg-1", B: "gfedcba-1", V: 1, T}`,
								`{B: "abcdefg-2", B: "gfedcba-2", V: 2, F}`,
							},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 1}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "DynamicBytes", Minor: 0},
						},
						Data: []string{
							`{B: "abcdefg-1"}`,
							`{B: "abcdefg-2"}`,
						},
					},
				},
				{
					Name:  "Cut all columns",
					Error: "nil", OpCode: "CUT",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								`{B: "abcdefg-1", B: "gfedcba-1", V: 1, T}`,
								`{B: "abcdefg-2", B: "gfedcba-2", V: 2, F}`,
							},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 0}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{},
						Data:  []string{"{}", "{}"},
					},
				},
				{
					Name:  "Cut error range - 1",
					Error: "errors.ErrorCodeIndexOutOfRange", OpCode: "CUT",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								`{B: "abcdefg-1", B: "gfedcba-1", V: 1, T}`,
								`{B: "abcdefg-2", B: "gfedcba-2", V: 2, F}`,
							},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 0}},
							Data:  []string{"{V: 5}"},
						},
					},
					Output: &tmplOp{},
				},
				{
					Name:  "Cut error range - 2",
					Error: "errors.ErrorCodeIndexOutOfRange", OpCode: "CUT",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								`{B: "abcdefg-1", B: "gfedcba-1", V: 1, T}`,
								`{B: "abcdefg-2", B: "gfedcba-2", V: 2, F}`,
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 0},
							},
							Data: []string{"{V: 15, V: 17}"},
						},
					},
					Output: &tmplOp{},
				},
			},
		},
		// -- end of CUT
		{
			TestName: "OpFilter", OpFunc: "opFilter",
			Cases: []*tmplTestCase{
				{
					Name:  "Filter first 2 rows",
					Error: "nil", OpCode: "FILTER",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
							Data:  []string{"{T, F}", "{F, T}", "{T, T}", "{F, F}"},
						},
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}},
							Data:  []string{"{T}", "{T}", "{F}", "{F}"},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{{Major: "Bool", Minor: 0}, {Major: "Bool", Minor: 0}},
						Data:  []string{"{T, F}", "{F, T}"},
					},
				},
			},
		},
		// -- end of FILTER
		{
			TestName: "OpCast", OpFunc: "opCast",
			Cases: []*tmplTestCase{
				{
					Name:  "None Immediate - int",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 1}, // int16 -> int8
								{Major: "Int", Minor: 1}, // int16 -> int24
							},
							Data: []string{"{V: 127, V: 127}", "{V: -128, V: -128}"},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
								{Major: "Int", Minor: 2},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 0},
							{Major: "Int", Minor: 2},
						},
						Data: []string{"{V: 127, V: 127}", "{V: -128, V: -128}"},
					},
				},
				{
					Name:  "None Immediate - int2",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 1}, // int16 -> uint16
								{Major: "Int", Minor: 1}, // int16 -> uint16
							},
							Data: []string{
								"{V: 32767, V: -32768}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 1},
								{Major: "Uint", Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Uint", Minor: 1},
							{Major: "Uint", Minor: 1},
						},
						Data: []string{
							"{V: 32767, V: 32768}",
						},
					},
				},
				{
					Name:  "None Immediate - int3",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 1}, // int16 -> bool
								{Major: "Int", Minor: 1}, // int16 -> bool
							},
							Data: []string{
								"{V: 32767, V: -32768}",
								"{V: 0, V: 0}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Bool", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Bool", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							"{T, T}", "{F, F}",
						},
					},
				},
				{
					Name:  "None Immediate - int4",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 1}, // int16 -> bytes16
								{Major: "Int", Minor: 1}, // int16 -> address
							},
							Data: []string{
								"{V: 32767, V: -32768}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "FixedBytes", Minor: 1},
								{Major: "Address", Minor: 0},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "FixedBytes", Minor: 1},
							{Major: "Address", Minor: 0},
						},
						Data: []string{
							"{B: {0x7f, 0xff}, B: {255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,0x80,0x00}}",
						},
					},
				},
				{
					Name:  "None Immediate - uint",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 1}, // uint16 -> uint8
								{Major: "Uint", Minor: 1}, // uint16 -> uint24
							},
							Data: []string{"{V: 128, V: 128}"},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 0},
								{Major: "Uint", Minor: 2},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Uint", Minor: 0},
							{Major: "Uint", Minor: 2},
						},
						Data: []string{"{V: 128, V: 128}"},
					},
				},
				{
					Name:  "None Immediate - uint2",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 1}, // uint16 -> int16
								{Major: "Uint", Minor: 1}, // uint16 -> byte16
							},
							Data: []string{
								"{V: 32767, V: 32768}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 1},
								{Major: "FixedBytes", Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 1},
							{Major: "FixedBytes", Minor: 1},
						},
						Data: []string{
							"{V: 32767, B: {0x80, 0x00}}",
						},
					},
				},
				{
					Name:  "None Immediate - uint3",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 1}, // uint16 -> bool
								{Major: "Uint", Minor: 1}, // uint16 -> bool
							},
							Data: []string{
								"{V: 32767, V: 0}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Bool", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Bool", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							"{T, F}",
						},
					},
				},
				{
					Name:  "None Immediate - uint4",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 1}, // uint16 -> bytes
								{Major: "Uint", Minor: 1}, // uint16 -> bytes
							},
							Data: []string{
								"{V: 32767, V: 0}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "FixedBytes", Minor: 1},
								{Major: "FixedBytes", Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "FixedBytes", Minor: 1},
							{Major: "FixedBytes", Minor: 1},
						},
						Data: []string{
							"{B: {0x7f, 0xff}, B: {0x00, 0x00}}",
						},
					},
				},
				{
					Name:  "None Immediate - uint5",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 1}, // uint16 -> address
							},
							Data: []string{
								"{V: 32767}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Address", Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Address", Minor: 1},
						},
						Data: []string{
							"{B: {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0x7f,0xff}}",
						},
					},
				},
				{
					Name:  "None Immediate - bytes",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "FixedBytes", Minor: 1}, // byte16 -> byte8
								{Major: "FixedBytes", Minor: 1}, // byte16 -> byte24
							},
							Data: []string{
								"{B: {0xff, 0xff}, B: {0xff, 0xff}}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "FixedBytes", Minor: 0},
								{Major: "FixedBytes", Minor: 2},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "FixedBytes", Minor: 0},
							{Major: "FixedBytes", Minor: 2},
						},
						Data: []string{
							"{B: {0xff}, B: {0xff, 0xff, 0x00}}",
						},
					},
				},
				{
					Name:  "None Immediate - bytes2",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "FixedBytes", Minor: 1}, // byte16 -> int16
								{Major: "FixedBytes", Minor: 1}, // byte16 -> uint16
							},
							Data: []string{
								"{B: {0x7f, 0xff}, B: {0x80, 0x00}}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 1},
								{Major: "Uint", Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 1},
							{Major: "Uint", Minor: 1},
						},
						Data: []string{
							"{V: 32767, V: 32768}",
						},
					},
				},
				{
					Name:  "None Immediate - bytes3",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "FixedBytes", Minor: 1}, // byte16 -> dyn
							},
							Data: []string{
								"{B: {0x7f, 0xff}}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "DynamicBytes", Minor: 1},
						},
						Data: []string{
							"{B: {0x7f, 0xff}}",
						},
					},
				},
				{
					Name:  "Same type",
					Error: "nil", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								"{T}",
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Bool", Minor: 0},
							},
							Data: []string{},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							"{T}",
						},
					},
				},
				{
					Name:  "Error Invalid Type",
					Error: "errors.ErrorCodeInvalidCastType", OpCode: "CAST",
					Inputs: []*tmplOp{
						{
							Im:    false,
							Metas: []*tmplOpMeta{{Major: "Int", Minor: 2}},
							Data:  []string{"{V: -32768}"},
						},
						{
							Im:    true,
							Metas: []*tmplOpMeta{{Major: "DynamicBytes", Minor: 0}},
							Data:  []string{},
						},
					},
					Output: &tmplOp{
						Im:    false,
						Metas: []*tmplOpMeta{},
						Data:  []string{},
					},
				},
			},
		},
		// -- end of CAST
		{
			TestName: "OpSort", OpFunc: "opSort",
			Cases: []*tmplTestCase{
				{
					Name:  "Multi-column sorting",
					Error: "nil", OpCode: "SORT",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								`{B: "c", V: 1, T}`,
								`{B: "b", V: 2, T}`,
								`{B: "a", V: 3, T}`,
								`{B: "a", V: 1, F}`,
								`{B: "b", V: 2, F}`,
								`{B: "c", V: 3, F}`,
								`{B: "b", V: 3, F}`,
								`{B: "a", V: 3, F}`,
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Bool", Minor: 0},
								{Major: "Uint", Minor: 1},
							},
							Data: []string{"{F, V: 1}", "{T, V: 2}", "{F, V: 0}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "DynamicBytes", Minor: 0},
							{Major: "Int", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							`{B: "c", V: 3, F}`,
							`{B: "b", V: 3, F}`,
							`{B: "a", V: 3, F}`,
							`{B: "a", V: 3, T}`,
							`{B: "b", V: 2, F}`,
							`{B: "b", V: 2, T}`,
							`{B: "a", V: 1, F}`,
							`{B: "c", V: 1, T}`,
						},
					},
				},
				{
					Name:  "Multi-column sorting - 2",
					Error: "nil", OpCode: "SORT",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "DynamicBytes", Minor: 0},
								{Major: "Int", Minor: 0},
								{Major: "Bool", Minor: 0},
							},
							Data: []string{
								`{B: "c", V: 1, T}`,
								`{B: "b", V: 2, T}`,
								`{B: "a", V: 3, T}`,
								`{B: "a", V: 1, F}`,
								`{B: "b", V: 2, F}`,
								`{B: "c", V: 3, F}`,
								`{B: "b", V: 3, F}`,
								`{B: "a", V: 3, F}`,
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Bool", Minor: 0},
								{Major: "Uint", Minor: 1},
							},
							Data: []string{"{T, V: 0}", "{T, V: 1}", "{F, V: 2}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "DynamicBytes", Minor: 0},
							{Major: "Int", Minor: 0},
							{Major: "Bool", Minor: 0},
						},
						Data: []string{
							`{B: "a", V: 1, F}`,
							`{B: "a", V: 3, T}`,
							`{B: "a", V: 3, F}`,
							`{B: "b", V: 2, T}`,
							`{B: "b", V: 2, F}`,
							`{B: "b", V: 3, F}`,
							`{B: "c", V: 1, T}`,
							`{B: "c", V: 3, F}`,
						},
					},
				},
			},
		},
		// -- end of SORT
		{
			TestName: "OpRange", OpFunc: "opRange",
			Cases: []*tmplTestCase{
				{
					Name:  "Range test limit 2 offset 1",
					Error: "nil", OpCode: "RANGE",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
							},
							Data: []string{
								`{V: 1}`, `{V: 2}`, `{V: 3}`, `{V: 4}`,
								`{V: 5}`, `{V: 6}`, `{V: 7}`, `{V: 8}`,
							},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 7},
								{Major: "Uint", Minor: 7},
							},
							Data: []string{"{V: 1, V: 2}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 0},
						},
						Data: []string{`{V: 2}`, `{V: 3}`},
					},
				},
				{
					Name:  "Range test limit 0 offset 1",
					Error: "nil", OpCode: "RANGE",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
							},
							Data: []string{`{V: 1}`, `{V: 2}`},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 7},
								{Major: "Uint", Minor: 7},
							},
							Data: []string{"{V: 1, V: 0}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 0},
						},
						Data: []string{},
					},
				},
				{
					Name:  "Range test offset 20",
					Error: "nil", OpCode: "RANGE",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
							},
							Data: []string{`{V: 1}`, `{V: 2}`},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 7},
							},
							Data: []string{"{V: 20}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 0},
						},
						Data: []string{},
					},
				},
				{
					Name:  "Range test limit 10 offset 20",
					Error: "nil", OpCode: "RANGE",
					Inputs: []*tmplOp{
						{
							Im: false,
							Metas: []*tmplOpMeta{
								{Major: "Int", Minor: 0},
							},
							Data: []string{`{V: 1}`, `{V: 2}`},
						},
						{
							Im: true,
							Metas: []*tmplOpMeta{
								{Major: "Uint", Minor: 7},
								{Major: "Uint", Minor: 7},
							},
							Data: []string{"{V: 20, V: 10}"},
						},
					},
					Output: &tmplOp{
						Im: false,
						Metas: []*tmplOpMeta{
							{Major: "Int", Minor: 0},
						},
						Data: []string{},
					},
				},
			},
		},
		// -- end of RANGE
	},
}
