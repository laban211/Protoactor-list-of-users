import csv
import random
import string
from random import randint

# def projNum_generator():
#   return ''.join(random.choice(string.digits) for _ in range(10))

projarr = []

def projNum_generator():
  i = 1
  while i < 15001:
    projarr.append(str(i).zfill(10))
    i += 1

  return projarr[int(random.uniform(0, 1) * 15000)]

def user_generator():
  userchars = string.ascii_lowercase + string.digits
  return ''.join(random.choice(userchars) for _ in range(8))

def week_generator():
  return str(random.randint(1,52))

def time_generator():
  return str(round(random.uniform(0, 9.9), 1))

def createRow():
  week = week_generator()
  user = user_generator()
  mon = time_generator()
  tue = time_generator()
  wed = time_generator()
  thu = time_generator()
  fri = time_generator()
  projNum = projNum_generator()
  
  row = (week, user, mon, tue, wed, thu, fri, projNum)
  return row

with open('text.csv', 'w', newline='') as csvfile:
    writer = csv.writer(csvfile, dialect='excel')
    for i in range(1000000):
      writer.writerow(createRow())