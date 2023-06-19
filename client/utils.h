#ifndef UTILS_H
#define UTILS_H

#include <string>

std::string login(std::string email, std::string password, std::string cookie_file);
std::string signUp(std::string email, std::string password);
std::string sendData(std::string data, std::string cookie_file);

#endif