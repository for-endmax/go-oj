result = ""
while True:
    line = input("")
    if line.lower() == 'exit':
        break
    result += line

print(result,end="")
