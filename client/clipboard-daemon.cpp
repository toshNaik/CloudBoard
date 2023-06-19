#include "utils.h"
#include <gtk/gtk.h>
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

std::string cookie_file;

void ClipboardTextChanged(GtkClipboard* clipboard, gpointer user_data)
{
	gchar* text = gtk_clipboard_wait_for_text(clipboard);
	// send json object to server
	std::string response = sendData(text, cookie_file);

	std::ofstream tmpfile;
	tmpfile.open("/tmp/cloudboard");
	tmpfile << response;
	g_free(text);
}


int main(int argc, char* argv[])
{
	// Login to server
	std::cout << "Enter email: ";
	std::string email;
	std::cin >> email;
	std::cout << "Enter password: ";
	std::string password;
	std::cin >> password;
	cookie_file = "/tmp/cookies.txt";
	// std::cin >> cookie_file;
	// if(cookie_file == "") {
	// 	cookie_file = "/tmp/cookies";
	// }
	
	std::string response = login(email, password, cookie_file);
	if (response[2] == 'e') {
		std::cout << response << std::endl;
		exit(0);
	} else {
		std::cout << "Logged in successfully" << std::endl;
	}

	// Daemonize
	pid_t pid, sid;
	pid = fork();
	if (pid < 0) {
		exit(EXIT_FAILURE);
	}
	if (pid > 0) {
		exit(EXIT_SUCCESS);
	}
	if (setsid() < 0) {
		exit(EXIT_FAILURE);
	}
	if ((chdir("/")) < 0) {
		exit(EXIT_FAILURE);
	}
	close(STDIN_FILENO);
	close(STDOUT_FILENO);
	close(STDERR_FILENO);

	gtk_init(&argc, &argv);

	GtkClipboard* clipboard = gtk_clipboard_get(GDK_SELECTION_CLIPBOARD);
	g_signal_connect(clipboard, "owner-change", G_CALLBACK(ClipboardTextChanged), NULL);

	gtk_main();

	return 0;
}

