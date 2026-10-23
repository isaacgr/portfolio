/**
 * Adds days to the current date
 *
 * @param {number} days
 * @returns {Date} new date with 'days' added to current date
 */
Date.prototype.addDays = function(days) {
  var date = new Date(this.valueOf());
  date.setDate(date.getDate() + days);
  return date;
};

/**
 * Get the total number of weeks between the date object and current date
 *
 * From:
 *  https://stackoverflow.com/questions/6117814/get-week-of-year-in-javascript-like-in-php
 *
 * Algorithm is to find nearest thursday, it's year
 * is the year of the week number. Then get weeks
 * between that date and the first day of that year.
 *
 * Note that dates in one year can be weeks of previous
 * or next year, overlap is up to 3 days.
 *
 * e.g. 2014/12/29 is Monday in week  1 of 2015
 *      2012/1/1   is Sunday in week 52 of 2011
 *
 * @see https://en.wikipedia.org/wiki/ISO_week_date
 * @see https://stackoverflow.com/questions/6117814/get-week-of-year-in-javascript-like-in-php
 *
 * @returns {number} Weeknumber for date object
 */
Date.prototype.getWeekNumber = function() {
  // Copy date so don't modify original
  const d = new Date(
    this.getUTCFullYear(),
    this.getUTCMonth(),
    this.getUTCDate()
  );
  // Set to nearest Thursday: current date + 4 - current day number
  // Make Sunday's day number 7
  d.setUTCDate(d.getUTCDate() + 4 - (d.getUTCDay() || 7));
  // Get first day of year
  var yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1));
  // Calculate full weeks to nearest Thursday
  return Math.ceil(((d - yearStart) / 86400000 + 1) / 7);
};

/**
 * Returns the number of weeks between two dates
 *
 * @param {Date} d1
 * @param {Date} d2
 * @returns {number} Number of weeks between d1 and d2
 */
function weeksBetween(d1, d2) {
  return Math.round((d2 - d1) / (7 * 24 * 60 * 60 * 1000));
}

function getDates(startDate, stopDate) {
  var dateArray = new Array();
  var currentDate = startDate;
  while (currentDate <= stopDate) {
    dateArray.push(new Date(currentDate));
    currentDate = currentDate.addDays(1);
  }
  return dateArray;
}

function isSameWeek(date1, date2) {
  const startDate = new LifeDate().startDate;
  if (weeksBetween(startDate, new Date(date1)) + 1 == date2) {
    return true;
  }
  return false;
}

function getWeekNumber(year, week) {
  return year * 52 - (53 - week) + 1;
}

function LifeDate() { }
LifeDate.prototype.startDate = new Date('1992-01-29');

LifeDate.prototype.datesArray = function() {
  return getDates(this.startDate, new Date());
};

LifeDate.prototype.born = function() {
  return {
    date: this.startDate,
    color: 'fill-white',
    name: 'Born'
  };
};

LifeDate.prototype.graduatedHighschool = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 17, 3, 16),
    color: 'fill-orange',
    name: 'Graduated Highschool'
  };
};

LifeDate.prototype.firstPostSecJob = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 24, 5, 16),
    color: 'fill-yellow',
    name: 'First Post Secondary Job'
  };
};

LifeDate.prototype.graduatedCollege = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 22, 3, 15),
    color: 'fill-red',
    name: 'Graduated College'
  };
};

LifeDate.prototype.graduatedUni = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 24, 3, 15),
    color: 'fill-blue',
    name: 'Graduated University'
  };
};

LifeDate.prototype.metWife = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 24, 5, 7),
    color: 'fill-pink',
    name: 'Met Wife'
  };
};

LifeDate.prototype.covid = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 28, 2, 17),
    color: 'fill-green',
    name: 'Covid'
  };
};

LifeDate.prototype.dadsDeath = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 62, 9, 29),
    color: 'fill-black',
    name: 'How Old Dad Was'
  };
};

LifeDate.prototype.dadsDeathDate = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 28, 9, 29),
    color: 'fill-black',
    name: 'Dads Death'
  };
};

LifeDate.prototype.momsDeathDate = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 33, 11, 20),
    color: 'fill-black',
    name: 'Moms Death'
  };
};

LifeDate.prototype.engaged = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 31, 5, 8),
    color: 'fill-purple',
    name: 'Engaged'
  };
};

LifeDate.prototype.firstHouse = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 31, 7, 8),
    color: 'fill-light-blue',
    name: 'First House'
  };
};

LifeDate.prototype.married = function() {
  return {
    date: new Date(this.startDate.getFullYear() + 32, 7, 8),
    color: 'fill-light-purple',
    name: 'Married'
  };
};

function getTdClassName(currYear, currWeek, numWeeks) {
  const lifeDate = new LifeDate();
  const gridWeek = getWeekNumber(currYear, currWeek);
  if (isSameWeek(lifeDate.dadsDeath().date, gridWeek)) return lifeDate.dadsDeath();
  else if (isSameWeek(lifeDate.dadsDeathDate().date, gridWeek)) return lifeDate.dadsDeathDate();
  else if (isSameWeek(lifeDate.momsDeathDate().date, gridWeek)) return lifeDate.momsDeathDate();
  else if (isSameWeek(lifeDate.metWife().date, gridWeek)) return lifeDate.metWife();
  else if (isSameWeek(lifeDate.graduatedCollege().date, gridWeek)) return lifeDate.graduatedCollege();
  else if (isSameWeek(lifeDate.graduatedUni().date, gridWeek)) return lifeDate.graduatedUni();
  else if (isSameWeek(lifeDate.born().date, gridWeek)) return lifeDate.born();
  else if (isSameWeek(lifeDate.covid().date, gridWeek)) return lifeDate.covid();
  else if (isSameWeek(lifeDate.graduatedHighschool().date, gridWeek)) return lifeDate.graduatedHighschool();
  else if (isSameWeek(lifeDate.firstPostSecJob().date, gridWeek)) return lifeDate.firstPostSecJob();
  else if (isSameWeek(lifeDate.engaged().date, gridWeek)) return lifeDate.engaged();
  else if (isSameWeek(lifeDate.firstHouse().date, gridWeek)) return lifeDate.firstHouse();
  else if (isSameWeek(lifeDate.married().date, gridWeek)) return lifeDate.married();
  else if (gridWeek <= numWeeks) return { color: 'fill', name: '' };
  else return { color: '', name: '' };
}

const maxAge = 85;

function renderAgeStats() {
  const age = new LifeDate().datesArray().length;
  document.getElementById('age').innerHTML = `Age: ${Math.floor(
    age / 7 / 52
  )} Years`;
  document.getElementById('time-lived').innerHTML = `
  Time Alive: ~${Math.floor(
    age / 7
  )} Weeks / ${age} Days`;
  document.getElementById('time-left').innerHTML = `
  Time Left: ~${Math.floor(maxAge - (age / 7) / 52)} Years / ${Math.floor(
    maxAge * 52 - age / 7
  )} Weeks / ${maxAge * 52 * 7 - age} Days`;
}

function generateGrid() {
  const datesArray = new LifeDate().datesArray();
  const numWeeks = Math.floor(datesArray.length / 7);
  let rows = [...Array(maxAge).keys()].map((_, year) => {
    let tds = [...Array(52).keys()]
      .map((_, week) => {
        let currYear = year + 1;
        let currWeek = week + 1;
        const { color: className, name: dateName } = getTdClassName(currYear, currWeek, numWeeks);
        const toolTip = `<span class="tooltiptext">${dateName}</span>`;
        return `<td id=${getWeekNumber(
          currYear,
          currWeek
        )} class="${className} tooltip">${(className.includes('fill-') && toolTip) || ''
          }</td>`;
      })
      .join('');
    return `<tr>${tds}</tr>`;
  });
  document.getElementById('life-grid').innerHTML = rows.join('');
}

generateGrid();
renderAgeStats();
