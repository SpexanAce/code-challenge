# Code Challenge

This repository contains an application (in multiple languages [Go, Python]) which calculates the average pollution for the morning, afternoon and night time, of PM10 and PM2.5 particle matter in Gothenburg for a specific date.

## The challenge

The challenge is multi faceted:

- The code has bugs. The average seems to be incorrect. The challenge is to correctly calculate the averages for the day.
- The application is currently very dumb. It requires a manual step to download the data for a specific date[1], and overwrite the database file in the repo. The challenge is to create a programmatic integration to [Open-Meteo](https://open-meteo.com/en/docs/air-quality-api) and use the previous days data.
- The application currently needs to be manually executed. The challenge is to expose the data via a REST API or a webpage.
- Try to think twice about the API design: how can the problem domain be extended? How can one incorporate other cities in the future? How can one incorporate other pollutants such as carbon monoxide in the future?

## The assignment

Complete the challenges in your language of choice. Try to have something completed in 4 hours time. That should be enough to understand the problem domain and to enable a meaningful discussion. It's to noones benefit to spend too much time on these challenges. If you still have unfinished challenges when the time is up, commit a note with status + how to complete the challenges (without actually completing them).

The assignment also includes committing your code changes to a git repository, as you go along (committing to `main`/`master` is fine). Create a zip file containing both the code and the `.git` folder when completing the assignment, and send in the zip file with e-mail.

---

[1] The URL which the data is downloaded from is https://air-quality-api.open-meteo.com/v1/air-quality?latitude=52.5235&longitude=13.4115&hourly=pm10,pm2_5&start_date=2023-01-31&end_date=2023-01-31