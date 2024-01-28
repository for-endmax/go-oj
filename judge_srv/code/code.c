#include <stdio.h>
#include <string.h>

#define MAX_LINE_LENGTH 1000

int main() {
    char line[MAX_LINE_LENGTH];
    char result[MAX_LINE_LENGTH * 10];

    while (1) {
        fgets(line, MAX_LINE_LENGTH, stdin);

        if (strcmp(line, "exit\n") == 0) {
            break;
        }
        int len = strlen(line);
        line[len-1]='\0';
        strcat(result, line);
    }

    printf("%s", result);

    return 0;
}