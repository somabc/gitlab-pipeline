# gitlab-pipeline-trigger

Trigger gitlab pipeline across projects using gitlab pipelines API, this tool makes
it easier to quick other projects pipeline. what you need to prepare:

1. Pipeline Trigger Token
2. API Token for your account/or deploy account
3. Project ID, that you want to trigger it's pipeline

The tool is contained in a small docker image tht you can use locally or inside gitlab ci:

`docker run --rm -it registry.gitlab.com/schehata/gitlab-pipeline-trigger ./main -a $API_TOKEN -t $TRIGGER_TOKEN -p $PROJECT_ID`

### Possible Attributes

 - `-a` to set the API Token
 - `-t` to set Trigger Token
 - `-h` to set host, default is `https://gitlab.com`
 - `-p` to set project id
 - `-b` to set target branch (reference), default is: master
 - `-w` to set if the client should keep fetching pipeline status until it the pipeline fails or succeeds,
disabled by default.

 ### Passing trigger variables

 From [gitlab docs](https://docs.gitlab.com/ce/ci/triggers/#making-use-of-trigger-variables):

    You can pass any number of arbitrary variables in the trigger API call and they will be available in GitLab CI so that they can be used in your .gitlab-ci.yml file.

This tools gives you the ability to pass these variables and converts them to `variables[key]=value` and send them to the
Gitlab pipeline api. you can pass those variables like this

```
docker run --rm -it registry.gitlab.com/schehata/gitlab-pipeline-trigger ./main env=production author=schehata logLevel=info
```

that would result in

```
variables[env]=production
variables[author]=schehata
variables[logLevel]=info
```

 ### Example of using it inside gitlab ci

 ```yml
e2e:
  stage: test
  when: manual
  image: registry.gitlab.com/ishehata/gitlab-pipeline-trigger
  script:
    - /trigger/main -a=$API_TOKEN -p=$PROJECT_ID -t=$TRIGGER_TOKEN -w=true
 ```
