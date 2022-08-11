package model

import (
	"testing"
)

func TestSanitizeLogUnique(t *testing.T) {
	testData := `Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
of this software and associated documentation files (the "Software"), to deal`
	want := `<line> permission is hereby granted <comma> free of charge <comma> to any person obtaining a copy </line> <line> of this software and associated documentation files <parenthesis> the <doublequote> software <doublequote> </parenthesis> <comma> to deal </line> <line> in the software without restriction <comma> including without limitation the rights </line>`
	result, err := SanitizeLog(testData, true)
	if want != result || err != nil {
		t.Fatalf(`SanitizeLog(...) = %q, %v, want match for %#q, nil`, result, err, want)
	}
}

func TestSanitizeLog(t *testing.T) {
	testData := `Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
python 3.8.10, pytest 1.2.3
2022/08/11 22:10:00 text
testing ... this type of line
testing...
Or a line of text.with a period in it
of this software and associated documentation files (the "Software"), to deal`
	want := `<line> permission is hereby granted <comma> free of charge <comma> to any person obtaining a copy </line> <line> of this software and associated documentation files <parenthesis> the <doublequote> software <doublequote> </parenthesis> <comma> to deal </line> <line> in the software without restriction <comma> including without limitation the rights </line> <line> python 3.8.10 <comma> pytest 1.2.3 </line> <line> <date> <time> text </line> <line> testing <period> <period> <period> this type of line </line> <line> testing <period> <period> <period> </line> <line> Or a line of text.with a period in it </line> <line> of this software and associated documentation files <parenthesis> the <doublequote> software <doublequote> </parenthesis> <comma> to deal </line>`
	result, err := SanitizeLog(testData, false)
	if want != result || err != nil {
		t.Fatalf(`SanitizeLog(...) = %q, %v, want match for %#q, nil`, result, err, want)
	}
}

func TestSpliceLines(t *testing.T) {
	testData := `<line> permission is hereby granted <comma> free of charge <comma> to any person obtaining a copy </line> <line> of this software and associated documentation files <parenthesis> the <doublequote> software <doublequote> </parenthesis> <comma> to deal </line> <line> in the software without restriction <comma> including without limitation the rights </line>`
	want := []string{
		"<line> permission is hereby granted <comma> free of charge <comma> to any person obtaining a copy </line>",
		"<line> of this software and associated documentation files <parenthesis> the <doublequote> software <doublequote> </parenthesis> <comma> to deal </line>",
		"<line> in the software without restriction <comma> including without limitation the rights </line>",
	}
	result := SpliceLines(testData)
	isSame := true
	if len(result) == len(want) {
		for i := 0; i < len(result); i++ {
			if result[i] != want[i] {
				isSame = false
			}
		}
	}
	if result == nil || len(result) != len(want) || !isSame {
		t.Fatalf(`SpliceLines(...) = %q, want match for %#q, nil`, result, want)
	}
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
// func TestHelloEmpty(t *testing.T) {
// 	msg, err := Hello("")
// 	if msg != "" || err == nil {
// 			t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
// 	}
// }