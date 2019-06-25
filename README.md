

## Dockerfiles Explained
 
There are three Dockerfiles represented in the Repo. If a change is being made to one Dockerfile, the other Dockerfiles should also be updated accordingly.

- Dockerfile 

    This Dockerfile is used by the Azure Pipeline to Build and Deploy the application

- Dockerfile_e2e_test:  
    
    This Dockerfile is used by the Azure Pipeline to run Integration level tests against the deployed application
    
- Dockerfile_local

    This Dockerfile is used to build a run a docker container on your local machine.  It does not contain the stage for uploacing test results to Azure
   
## Develop and Run Slackbot

To begin developing against this repository, the following steps should be followed.

1) Setup environment variables    

        export GOPATH=<Path on your machine where go modules will be stored>
        
        Obtain the two tokens that are passed into the app are used to verify API requests.  Both tokens are registered with Slack for the app and will be sent with API requests. However, each serve a different purpose. The Verification token is used to validate Domain changes in slack. When the event url is changed [here](https://api.slack.com/apps/AG29FUH1U/event-subscriptions?), the bot will handle the request and echo the token back to validate communication.  The Oauth token is used by the app to validate all other API requests coming from Slack. 
            export APPSETTING_SLACKBOT_OAUTHTOKEN: <value below>
            export APPSETTING_SLACKBOT_VERIFICATIONTOKEN: <value below>
    
            APSETTING_SLACKBOT_VERIFICATIONTOKEN: https://api.slack.com/apps/AG29FUH1U/general?  
            APPSETTING_SLACKBOT_OAUTHTOKEN: https://api.slack.com/apps/AG29FUH1U/oauth?
    
2) Clone to Repo
 
        git clone https://github.com/buildit/slackbot.git

3) Build the Container    

        cd <workspace>/slackbot  
        
        docker build --build-arg APPSETTING_SLACKBOT_OAUTHTOKEN=${APPSETTING_SLACKBOT_OAUTHTOKEN} --build-arg APPSETTING_SLACKBOT_VERIFICATIONTOKEN=${APPSETTING_SLACKBOT_VERIFICATIONTOKEN} --target=final -f Dockerfile_local -t slackbot:latest .
        
        
4) Run the Container

        docker run -d -p 4390:4390 -e APPSETTING_SLACKBOT_OAUTHTOKEN=${APPSETTING_SLACKBOT_OAUTHTOKEN} -e APPSETTING_SLACKBOT_VERIFICATIONTOKEN=${APPSETTING_SLACKBOT_VERIFICATIONTOKEN} slackbot:latest  
        
5) Configure slack to communicate with your slackbot container that is now running on localhost

        Once the app is running, Slack needs to be pointed at the running container. For running the app locally on your machine you can establish a tunnel to a port on your machine using NGROK. However, when you point slack at this local domain (step 4 below), the hosted app will no longer be receiving API events.  Note: This works for the time being, but if the bot/app gets higher usage, this local development and repointing of slack will not suffice.
        1) Download NGROK to enable local development: https://ngrok.com/download  
        2) Install it to /usr/local/bin or somewhere on your path  
        3) Slackbot runs on port 4390, so you'd run ngrok in a terminal to expose that port on your machine ex. ngrok http 4390  
        4) Grab the forwarding domain (ex. http://933d7ddb.ngrok.io) and use that as the domain name for the three different configurations in slack for events, interactions, and slash commands:
         - events:  https://api.slack.com/apps/AG29FUH1U/event-subscriptions?  
         - interactions:  https://api.slack.com/apps/AG29FUH1U/interactive-messages?  
         - slash commands:  https://api.slack.com/apps/AG29FUH1U/slash-commands?  
         
 6) Invite @miles to a slack channel.  Type "@miles hi" and the bot should respond "Hello" in your channel.
 
Note:  If you would rather run main.go in your IDE instead of continuing to Build and Run the container, you just need to ensure that your Run Configuration sets the environmnet variables:
        
        APPSETTING_SLACKBOT_OAUTHTOKEN
        APPSETTING_SLACKBOT_VERFICATIONTOKEN
        
TODO: 

1)Add a help menu as we expand bot functionality and event processing. Currently if you mention the app (ex. @miles), the bot will just respond with "Hello"
