#include "utils.h"
#include <iostream>
#include <fstream>
#include <sys/types.h>
#include <sys/stat.h>
#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <errno.h>
#include <unistd.h>
#include <syslog.h>
#include <curl/curl.h>

// Callback function to write response data into a string
size_t writeCallback(void* contents, size_t size, size_t nmemb, std::string* response) {
    size_t totalSize = size * nmemb;
    response->append((char*)contents, totalSize);
    return totalSize;
}

std::string sendData(std::string data, std::string cookie_file)
{
    CURL* curl = curl_easy_init();
    std::string response_string;
    if (curl) {
        std::string url = "https://cloudboard-389721.uc.r.appspot.com/cloudboard/put";
        std::string jsonData = "{\"data\":\"" + data + "\"}";

        // set request headers
        struct curl_slist* headers = NULL;
        headers = curl_slist_append(headers, "Content-Type: application/json");

        // set curl options
        curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
        curl_easy_setopt(curl, CURLOPT_POST, 1L);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, jsonData.c_str());
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, writeCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response_string);
        curl_easy_setopt(curl, CURLOPT_COOKIEFILE, cookie_file.c_str());

        // perform curl request
        CURLcode res = curl_easy_perform(curl);
        if (res != CURLE_OK) {
            std::cerr << "curl_easy_perform() failed: " << curl_easy_strerror(res) << std::endl;
        }

        // cleanup
        curl_easy_cleanup(curl);
        curl_slist_free_all(headers);
    }
    return response_string;
}

std::string signUp(std::string email, std::string password)
{
	CURL* curl = curl_easy_init();
    std::string response_string;
	if (curl) {
		std::string url = "https://cloudboard-389721.uc.r.appspot.com/signup";
		std::string jsonData = "{\"email\":\"" + email + "\",\"password\":\"" + password + "\"}";
  
		// set request headers
		struct curl_slist* headers = NULL;
		headers = curl_slist_append(headers, "Content-Type: application/json");

		// set curl options
		curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
		curl_easy_setopt(curl, CURLOPT_POST, 1L);
		curl_easy_setopt(curl, CURLOPT_POSTFIELDS, jsonData.c_str());
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, writeCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response_string);
        curl_easy_setopt(curl, CURLOPT_COOKIEJAR, "cookies.txt");

        // perform curl request
        CURLcode res = curl_easy_perform(curl);

        if (res != CURLE_OK) {
            std::cerr << "curl_easy_perform() failed: " << curl_easy_strerror(res) << std::endl;
            return "";
        }
        // cleanup
		curl_easy_cleanup(curl);
        curl_slist_free_all(headers);
	}
    return response_string;
}


std::string login(std::string email, std::string password, std::string cookie_file)
{
	CURL* curl = curl_easy_init();
    std::string response_string;
	if (curl) {
		std::string url = "https://cloudboard-389721.uc.r.appspot.com/login";
		std::string jsonData = "{\"email\":\"" + email + "\",\"password\":\"" + password + "\"}";
  
		// set request headers
		struct curl_slist* headers = NULL;
		headers = curl_slist_append(headers, "Content-Type: application/json");

		// set curl options
		curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
		curl_easy_setopt(curl, CURLOPT_POST, 1L);
		curl_easy_setopt(curl, CURLOPT_POSTFIELDS, jsonData.c_str());
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, writeCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response_string);
        curl_easy_setopt(curl, CURLOPT_COOKIEJAR, cookie_file.c_str());

        // perform curl request
        CURLcode res = curl_easy_perform(curl);

        if (res != CURLE_OK) {
            std::cerr << "curl_easy_perform() failed: " << curl_easy_strerror(res) << std::endl;
            return "";
        }
        // cleanup
		curl_easy_cleanup(curl);
        curl_slist_free_all(headers);
	}
    return response_string;
}
