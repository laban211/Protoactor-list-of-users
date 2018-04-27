import csv
import random
import string
from random import randint

projectNumbers = []
users = [] 

def projNum_generator():
  i = 1
  while i < 15001:
    projectNumbers.append(str(i).zfill(10))
    i += 1

  return projectNumbers

def getRandomProjNum(arr):
  return arr[int(random.uniform(0,1)*15000)]

def user_generator():
  userchars = string.ascii_lowercase + string.digits
  i = 1
  while i < 3001:
    users.append("user" + str(i))
    i += 1
  return users

def getRandomUser(arr):
  return arr[int(random.uniform(0,1)*3000)]

def week_generator():
  return str(random.randint(1,52))

def time_generator():
  return str(round(random.uniform(1, 8), 1))

def createRow():
  week = week_generator()
  user = getRandomUser(theUsers)
  mon = time_generator()
  tue = time_generator()
  wed = time_generator()
  thu = time_generator()
  fri = time_generator()
  projNum = getRandomProjNum(theProjnums)
  
  row = (week, user, mon, tue, wed, thu, fri, projNum)
  return row

theProjnums = projNum_generator()
theUsers = user_generator()



with open('text.csv', 'w', newline='') as csvfile:
    writer = csv.writer(csvfile, dialect='excel')
    for i in range(1000000):
      writer.writerow(createRow())