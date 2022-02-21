/*MIT License

Copyright (©) 2019 - Randall Simpson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.*/
package model

import (
  "strings"
  "regexp"
  "fmt"
  "sort"
)

func SanitizeText(data string) (string, error) {
  removePunctuation := func(r rune) rune {
		if strings.ContainsRune(".,:;", r) {
			return -1
		} else {
			return r
		}
	}
  
  s := strings.Map(removePunctuation, data)
  
  return s, nil
}

func SanitizeLog(data string, unique bool) (string, error) {
	//replacer := strings.NewReplacer(",", " <comma> ", ".", " </s> ", ";", " <colon> ")
  //replacer := strings.NewReplacer("\n", " </line> ")
  s := data
  
  re := regexp.MustCompile(`([E][0-9]{4})`)
  s = re.ReplaceAllString(s, " <error> ")
  re = regexp.MustCompile(`([I][0-9]{4})`)
  s = re.ReplaceAllString(s, " <info> ")
  //do date-utc
  re = regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2} ([0-9]{2}:)+[0-9]{2}.[0-9]* \+[0-9]* \w*`)
  s = re.ReplaceAllString(s, " <date-utc> ")
  //end of file </file> \w*.\w*:[0-9]*\] and pid
  re = regexp.MustCompile(` *(\d*) (\w*.\w*:[0-9]*)\]`)
  s = re.ReplaceAllString(s, " <pid> ${1} </pid> <file> ${2} </file> ")
  //dates in the beginning of the line
  re = regexp.MustCompile(`([0-9]{2}:)+[0-9]{2}.[0-9]*`)
  s = re.ReplaceAllString(s, " <date> ")
  
  //find numbers except for within text or the file locations.
  re = regexp.MustCompile(`[ [+]\d+[.]*\d+[\]ms]*`)
  s = re.ReplaceAllString(s, " <number> ")

  re = regexp.MustCompile(`[\t]`)
  s = re.ReplaceAllString(s, " <tab> ")

  re = regexp.MustCompile(`[,]`)
  s = re.ReplaceAllString(s, " <comma> ")

  re = regexp.MustCompile(`[.] +`)
  s = re.ReplaceAllString(s, " <period> ")

  re = regexp.MustCompile(`[;]`)
  s = re.ReplaceAllString(s, " <semicolin> ")

  re = regexp.MustCompile(`["]`)
  s = re.ReplaceAllString(s, " <doublequote> ")

  re = regexp.MustCompile(`[(]`)
  s = re.ReplaceAllString(s, " <parenthesis> ")

  re = regexp.MustCompile(`[)]`)
  s = re.ReplaceAllString(s, " </parenthesis> ")

  //find \r and replace them
  re = regexp.MustCompile(`[\r\n]`)
  s = re.ReplaceAllString(s, "\n")
  
  //double spaces or more into single space.
  re = regexp.MustCompile(` {2,}`)
  s = re.ReplaceAllString(s, " ")
  
  //lowercase
  s = strings.ToLower(s)
  
  //unique string for each line.
  set := make(map[string]int)
  
  lines := strings.Split(s, "\n")

  s = ""
	for i, line := range lines {
    lineOfText := strings.TrimSpace(line)
    if len(lineOfText) > 0 {
      if unique {
        _, exists := set[lineOfText]
        if !exists {
          set[lineOfText] = i
        }
      } else {
        s += fmt.Sprintf("<line> %s </line> ", lineOfText)
      }
    }
	}

  if unique {
    //ensure the output is in order (less the duplicates)
    keys := make([]string, 0, len(set))
    for key := range set {
      keys = append(keys, key)
    }
    sort.Slice(keys, func(i, j int) bool { return set[keys[i]] < set[keys[j]] })

    for _, key := range keys {
      s += fmt.Sprintf("<line> %s </line> ", key)
    }
  }

  //remove the last space
  if len(s) > 0 {
    s = s[:len(s) - 1]
  }
  
  return s, nil
}

func SpliceLines(data string) []string {
  replacer := strings.NewReplacer("</line> <line>", "</line>\n<line>")
	s := replacer.Replace(data)

  return strings.Split(s, "\n")
}