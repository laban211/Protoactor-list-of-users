import csv

found = {}

with open('text.csv', newline='') as csvfile:
  reader = csv.reader(csvfile, dialect='excel')
  for row in reader:
    # splitted = row.split()
    projNo = row[7]
    if projNo not in found:
      found[projNo] = 1
    else:
      found[projNo] = found[projNo] + 1

sortedfound = sorted(found, key = found.get, reverse=True)


topscore = sortedfound[0]
# print(topscore)
for i in sortedfound:
  if found[i] == found[topscore]:
    print(i, found[i])
  else:
    break


# dict
# varje rad skapa key som är projnunmmer

#for
# if()counter++ 
# när vi hittar projnr