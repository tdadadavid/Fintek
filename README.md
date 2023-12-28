# PhynTek

****
## SETUP
****

### Clone project

```git
  🍕. git clone git@github.com:tdadadavid/Fintek.git
```

### Install dependencies
```markdown
  🍕. make download_deps
```

### Create environment files and file them correctly
```bash
  🍕. cp .env.env.example env.env
```

### Start project [dev mode]
```bash
  🍕. make start_dev
```

### Start project [prod]
```bash
  🍕. make start_prod
```

### Using docker 🐳
```bash
  🐳 docker run dockerrundavid/fingreat:latest
```

**** 
# GUIDELINES
****
* When creating *new branch* if you're fixing a bug or implementing a feature follow this pattern
  - fixing: *`fix/<name-of-fix>`*
  - feature: *`feat/<name-of-feature>`*
* Never commit directly into the *`main`* branch.



