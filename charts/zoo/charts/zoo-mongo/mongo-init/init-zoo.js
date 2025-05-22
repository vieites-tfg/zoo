use("zoo");

db.createCollection("animals");

db.animals.insertMany([
  {
    name: "Leo",
    species: "Lion",
    birthday: "2018-03-15",
    genre: "male",
    diet: "Carnivore",
    condition: "Healthy",
    notes: "Alpha male in the pride."
  },
  {
    name: "Nala",
    species: "Lion",
    birthday: "2019-07-10",
    genre: "female",
    diet: "Carnivore",
    condition: "Healthy",
    notes: "Sister of the alpha male."
  },
  {
    name: "Zuri",
    species: "Giraffe",
    birthday: "2017-11-01",
    genre: "female",
    diet: "Herbivore",
    condition: "Healthy"
  },
  {
    name: "George",
    species: "Giraffe",
    birthday: "2016-06-22",
    genre: "male",
    diet: "Herbivore",
    condition: "Healthy"
  },
  {
    name: "Tisha",
    species: "Elephant",
    birthday: "2012-12-30",
    genre: "female",
    diet: "Herbivore",
    condition: "Healthy"
  },
  {
    name: "Dumbo",
    species: "Elephant",
    birthday: "2010-05-12",
    genre: "male",
    diet: "Herbivore",
    condition: "Injured",
    notes: "Recovering from minor foot injury."
  },
  {
    name: "Stripes",
    species: "Zebra",
    birthday: "2018-01-19",
    genre: "male",
    diet: "Herbivore",
    condition: "Healthy"
  },
  {
    name: "Zara",
    species: "Zebra",
    birthday: "2019-02-08",
    genre: "female",
    diet: "Herbivore",
    condition: "Healthy"
  },
  {
    name: "Miko",
    species: "Monkey",
    birthday: "2020-09-15",
    genre: "male",
    diet: "Omnivore",
    condition: "Healthy",
    notes: "Very playful with visitors."
  },
  {
    name: "Kali",
    species: "Monkey",
    birthday: "2021-03-03",
    genre: "female",
    diet: "Omnivore",
    condition: "Healthy"
  },
  {
    name: "Ray",
    species: "Parrot",
    birthday: "2019-04-10",
    genre: "male",
    diet: "Omnivore",
    condition: "Healthy",
    notes: "Knows a few words."
  },
  {
    name: "Rita",
    species: "Parrot",
    birthday: "2018-11-18",
    genre: "female",
    diet: "Omnivore",
    condition: "Healthy"
  },
  {
    name: "Cleo",
    species: "Tiger",
    birthday: "2015-07-02",
    genre: "female",
    diet: "Carnivore",
    condition: "Healthy"
  },
  {
    name: "Sheru",
    species: "Tiger",
    birthday: "2014-09-29",
    genre: "male",
    diet: "Carnivore",
    condition: "Healthy",
    notes: "Recently transferred from another zoo."
  },
  {
    name: "Benny",
    species: "Penguin",
    birthday: "2019-12-25",
    genre: "male",
    diet: "Carnivore",
    condition: "Healthy"
  },
  {
    name: "Penny",
    species: "Penguin",
    birthday: "2020-12-01",
    genre: "female",
    diet: "Carnivore",
    condition: "Healthy"
  },
  {
    name: "Lola",
    species: "Flamingo",
    birthday: "2017-08-08",
    genre: "female",
    diet: "Omnivore",
    condition: "Healthy"
  },
  {
    name: "Felipe",
    species: "Flamingo",
    birthday: "2016-02-14",
    genre: "male",
    diet: "Omnivore",
    condition: "Healthy"
  },
  {
    name: "Gaston",
    species: "Hippo",
    birthday: "2011-03-23",
    genre: "male",
    diet: "Herbivore",
    condition: "Ill",
    notes: "Under veterinary supervision."
  },
  {
    name: "Hilda",
    species: "Hippo",
    birthday: "2011-08-30",
    genre: "female",
    diet: "Herbivore",
    condition: "Healthy"
  }
]);

print("'ZOO' DATABASE INITIALIZED WITH 20 ANIMALS IN 'ANIMALS' COLLECTION.");
