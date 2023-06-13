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

void ClipboardTextChanged(GtkClipboard* clipboard, gpointer user_data)
{
	gchar* text = gtk_clipboard_wait_for_text(clipboard);
	std::ofstream tmpfile;
	tmpfile.open("/tmp/cloudboard");
	// g_print("Clipboard contents: %s\n", text);
	tmpfile << text;
	g_free(text);
}

int main(int argc, char* argv[])
{
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

