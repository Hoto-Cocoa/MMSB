## MMSB
Misskey Mention Spamming Blocker

## Usage
1. Download the latest release from [here](https://github.com/Hoto-Cocoa/MMSB/releases).
2. Run the `run.bat` file.
3. Follow the instructions.

## How to create token?
1. Go to [API Tokens (You can select your instance at link)](https://misskey-hub.net/mi-web/?path=/settings/api).
2. Click "API Console".
3. Use `miauth/gen-token` for Endpoint, and paste this code in params and click "Send".
```javascript
{
  session: null,
  name: 'MMSB',
  description: 'Misskey Mention Spamming Blocker',
  permission: ['read:account', 'write:notes', 'write:admin:suspend-user'],
}
```
3. Copy the token and paste it in the program.

## License
See [LICENSE](/LICENSE) file for details.

## Special Thanks
- [nulta](https://github.com/nulta): Created the original version of this program.
