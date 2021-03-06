package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func print404(w *http.ResponseWriter, customMessage interface{}) {
	log.Println(customMessage)
	(*w).WriteHeader(http.StatusNotFound)
	(*w).Write([]byte("404 page not found"))
}

func handleChannelRequest(w http.ResponseWriter, r *http.Request) {
	reqPath := strings.Replace(r.URL.Path, "/iptv/", "", 1)
	reqPathParts := strings.SplitN(reqPath, "/", 2)
	reqPathPartsLen := len(reqPathParts)

	// Exit if no channel and/or no path provided:
	if reqPathPartsLen == 0 {
		print404(&w, "Unable to properly extract data from request '"+r.URL.Path+"'!")
		return
	}

	// Remove ".m3u8" from channel name
	if reqPathPartsLen == 1 {
		reqPathParts[0] = strings.Replace(reqPathParts[0], ".m3u8", "", 1)
	}

	// Extract channel name:
	encodedChannelName := &reqPathParts[0]
	decodedChannelName, err := url.QueryUnescape(*encodedChannelName)
	if err != nil {
		print404(&w, "Unable to decode channel '"+*encodedChannelName+"'!")
		return
	}

	// Retrieve channel from channels map:
	tvMutex.RLock()
	channel, ok := tvChannels[decodedChannelName]
	tvMutex.RUnlock()
	if !ok {
		print404(&w, "Unable to find channel '"+decodedChannelName+"'!")
		return
	}

	// For channel we need URL. For anything else we need URL root:
	var requiredURL string
	tvMutex.RLock()
	if reqPathPartsLen == 1 {
		requiredURL = channel.URL
	} else {
		requiredURL = channel.URLRoot + reqPathParts[1]
	}
	tvMutex.RUnlock()
	if requiredURL == "" {
		print404(&w, "Channel '"+decodedChannelName+"' does not have URL assigned!")
		return
	}

	// Retrieve requiredURL contents
	resp, err := http.Get(requiredURL)
	if err != nil {
		print404(&w, err)
		return
	}
	defer resp.Body.Close()

	// If not code 200
	if resp.StatusCode != 200 {
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte("Error"))
		return
	}

	// If path ends with ".ts" - return raw fetched bytes
	if strings.HasSuffix(r.URL.Path, ".ts") {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			print404(&w, err)
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}

	// Write everything, but rewrite links to itself
	w.WriteHeader(resp.StatusCode)
	prefix := "http://" + r.Host + "/iptv/" + *encodedChannelName + "/"
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			line = prefix + line
		} else if strings.Contains(line, "URI=\"") && !strings.Contains(line, "URI=\"\"") {
			line = strings.ReplaceAll(line, "URI=\"", "URI=\""+prefix)
		}
		w.Write([]byte(line + "\n"))
	}
}
