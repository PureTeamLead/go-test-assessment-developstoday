# Spy Cat Agency Management Service

Requirements:
- Golang SDK
- Bash 
- Docker (Docker Desktop running)
- GNU utils such as **make**

## Run the application

- Clone the repository
```bash
mkdir spy_cat_agency
cd spy_cat_agency
git clone https://github.com/PureTeamLead/go-test-assessment-developstoday .
 ```
- Create .env file from .env.example template(for better debugging purposes leave db name).
- Run the app
```bash
docker-compose up -d
```

## Possible issues

- Error with migrations
Solution:
Let Postgres container run and open another Terminal window:

```bash
make mig-down
```

Then run again:
```bash
docker-compose up -d
```


## Application endpoints

Are available as a [Postman Collection](https://www.postman.com/payload-candidate-40552263/spycatagency/collection/le9jv1k/spycatagency?action=share&creator=40502373)
