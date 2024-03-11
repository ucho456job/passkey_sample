import { encode } from 'base64url';

function base64UrlStringToArrayBuffer(base64UrlString) {
  let base64 = base64UrlString.replace(/-/g, '+').replace(/_/g, '/');
  let binaryString = window.atob(base64);
  let bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes.buffer;
}

function arrayBufferToBase64Url(buffer) {
  var binary = '';
  var bytes = new Uint8Array(buffer);
  var len = bytes.byteLength;
  for (var i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  var base64 = window.btoa(binary);
  return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

export default function Register() {
  const handleGeneratePasskey = async () => {
    console.log('handleGeneratePasskey')
    try {
      const challengeResponse = await fetch('/api/auth/challenge', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
      });
      const challengeData = await challengeResponse.json();
      console.log('challenge response: ', challengeData);

      challengeData.publicKey.challenge = base64UrlStringToArrayBuffer(challengeData.publicKey.challenge);
      challengeData.publicKey.user.id = base64UrlStringToArrayBuffer(challengeData.publicKey.user.id);

      console.log('encoded challenge response: ', challengeData);

      const result = await navigator.credentials.create(challengeData);
      console.log("navigatore.credentials.create result: ", result);

      const credentials = {
        id: result.id,
        type: result.type,
        rawId: arrayBufferToBase64Url(result.rawId),
        response: {
          clientDataJSON: arrayBufferToBase64Url(result.response.clientDataJSON),
          attestationObject: arrayBufferToBase64Url(result.response.attestationObject),
        },
      };
      console.log('create passkey request body: ', credentials);

      const passkeyResponse = await fetch('/api/auth/passkey', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      });

      const passkeyData = await passkeyResponse.json();
      console.log('create passkey response: ', passkeyData);
      alert('Passkey created successfully!')
    } catch (err) {
      console.error('Error:', err);
      alert('An error occurred.');
    }
  };

  return (
    <div>
      <h1>Register</h1>
      <button onClick={handleGeneratePasskey}>Generate Passkey</button>
    </div>
  );
}
