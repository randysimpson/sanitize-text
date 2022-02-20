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
package api
 
import (
    "encoding/json"
    "fmt"
    "net/http"
    "sanitize-text/model"
    "io/ioutil"
    "io"
    "k8s.io/klog"
    "strings"
    "bytes"
)

type SanitizedText struct {
	Data string		`json:"data"`
	Size int 	    `json:"size"`
	Status string `json:"status"`
}

type SanitizedTextArray struct {
  Lines []string `json:"lines"`
  Size int 	    `json:"size"`
	Status string `json:"status"`
}
 
func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Welcome!")
}
 
func LoadText(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    klog.Errorf("error: %+v", err)
  }
  
  if err := r.Body.Close(); err != nil {
    klog.Errorln(err)
  }

  data, err := model.SanitizeText(string(body))
  if err != nil {
    klog.Errorf("error: %+v", err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  
  response := SanitizedText{data, len(data), "Success"}

  jsonResponse, err := json.Marshal(response)
  if err != nil {
    klog.Errorf("error: %+v", err)
    return
  }

  w.Write(jsonResponse)
}

func LoadLogs(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    klog.Errorf("error: %+v", err)
  }
  
  if err := r.Body.Close(); err != nil {
    klog.Errorln(err)
  }

  data, err := model.SanitizeLog(string(body))
  if err != nil {
    klog.Errorf("error: %+v", err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  
  response := SanitizedText{data, len(data), "Success"}

  jsonResponse, err := json.Marshal(response)
  if err != nil {
    klog.Errorf("error: %+v", err)
    return
  }

  w.Write(jsonResponse)
}

func UploadLogs(w http.ResponseWriter, r *http.Request) {
  r.ParseMultipartForm(32 << 20) // limit your max input length!
  var buf bytes.Buffer
  
  file, header, err := r.FormFile("file")
  if err != nil {
    klog.Errorf("error: %+v", err)
  }
  defer file.Close()
  name := strings.Split(header.Filename, ".")
  klog.Infof("File name %s\n", name[0])
  
  // Copy the file data to my buffer
  io.Copy(&buf, file)
  
  contents := buf.String()
  data, err := model.SanitizeLog(contents)
  if err != nil {
    klog.Errorf("error: %+v", err)
  }
  
  // I reset the buffer in case I want to use it again
  // reduces memory allocations in more intense projects
  buf.Reset()
  
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)

  keys, ok := r.URL.Query()["lines"]
  
  if !ok || len(keys[0]) < 1 || keys[0] != "true" {
    response := SanitizedText{data, len(data), "Success"}

    jsonResponse, err := json.Marshal(response)
    if err != nil {
      klog.Errorf("error: %+v", err)
      return
    }

    w.Write(jsonResponse)
  } else {
    response := SanitizedTextArray{model.SpliceLines(data), len(data), "Success"}

    jsonResponse, err := json.Marshal(response)
    if err != nil {
      klog.Errorf("error: %+v", err)
      return
    }

    w.Write(jsonResponse)
  }
}