#include "utils.h"
#include <iostream>

int main(int argc, char* argv[])
{
    // Signup and login to server
    std::cout << "Enter email: ";
    std::string email;
    std::cin >> email;
    std::cout << "Enter password: ";
    std::string password;
    std::cin >> password;
    // curl email and password to server as json
    std::string response = signUp(email, password);
    std::cout << response << std::endl;
    return 0;
}

