# GoGoDrive

**gogodrive** is a very simple CLI application to upload files to your Google Drive account.

It's strongly based on this [tutorial](https://developers.google.com/drive/api/v3/quickstart/go) (which I recommend in order to know how to create 
the proper credentials to allow the application access you GDrive account) and on this [post](https://medium.com/@devtud/upload-files-in-google-drive-with-golang-and-google-drive-api-d686fb62f884). Basically, I just glue all together.

## **Why??**

Just for fun. And I needed a tool to upload files to a GDrive account, periodically (via cron) on different Linux machines with different Linux flavours.
And I want to learn Golang, so this tool seemed to be a good start. And, you know, just for fun.

## How to use

Assuming that you already have the credential file, you also need a configuration file (in JSON) like:

```javascript
{
	"credential.file": "absolute path to google credendial file JSON",
	"token.file": "absolute path to token file",
	"delete.pattern": "search file string",
	"upload.max_tries": 3 // max number of attempts to upload the file, if fail
}
```

This tool was designed to upload fresh backups on a GDrive account, trying not to exceed the account space (unless the file you trying to upload is greater than the left space in your GDrive), so the **delete.pattern** is used to **"delete"** older files before upload a new one.

**delete.pattern** follow the Google Drive API query patterns described [here](https://developers.google.com/drive/api/v3/reference/query-ref), with a plus, you can add a variable *today* date in the pattern:
```javascript
{
    // ...
    "delete.pattern": "modifiedTime <= '{today-10}'"
}
```
That will delete all files where *modifiedTime* is less or equal than 10 days before the machine current date.

Sure, you can leave **delete.pattern** empty and never delete files at all.

**Usage**

```bash
gogodrive -c config_file -i upload.file [-o target_filename]
```
If the argument *-o* is missing, the target file name will be the same as uploaded file name.

If you don't have the *token.file* yet, the application will ask you to go to a link in your browser in order to authorize access to your GDrive account. Copy the authorization code and paste it in the terminal, then press enter.

```
Go to the following link in your browser then type the authorization code: 
https://accounts.google.com/.....
[paste the authorization code here and press enter]
```

After that, the token file will be created and saved in the path you set on *token.file* property of your config.json file. Once created, the token file won't be asked anymore. So, before put *gogodrive* in a crontab, you will need:

- Go to [https://console.cloud.google.com](https://console.cloud.google.com) in order to create the *credential* file (see [here](https://developers.google.com/drive/api/v3/quickstart/go));

- Run *gogodrive* once by yourself in order to create the *token* file;

- With these two files, you can put *gogodrive* to run in a crontab.

## Tests

Yeah! I know. Sorry.
