require('dotenv').config();  // Load environment variables from .env file

const express = require('express');
const bodyParser = require('body-parser');
const crypto = require('crypto');

const app = express();
const PORT = process.env.PORT || 3000;

// Body parser middleware to handle raw JSON data
app.use(bodyParser.json({
    verify: (req, res, buf) => {
        req.rawBody = buf;
    }
}));

// Middleware to verify HubSpot signature
app.use('/hubspot-webhook', (req, res, next) => {
    const signature = req.headers['x-hubspot-signature'];
    const appSecret = process.env.CLIENT_SECRET; // Use CLIENT_SECRET from your .env file
    const httpMethod = req.method.toUpperCase();
    const requestUri = req.originalUrl;
    const requestBody = req.rawBody.toString();

    // Concatenate app secret, HTTP method, URI, and request body to form the base string
    const baseString = appSecret + httpMethod + requestUri + requestBody;

    // Create a SHA-256 hash of the base string
    const hash = crypto.createHash('sha256').update(baseString).digest('hex');

    // Compare the computed hash with the provided signature
    if (hash === signature) {
        console.log('Signature verified');
        next();
    } else {
        console.log('Invalid signature');
        res.status(403).send('Invalid signature');
    }
});

// Main handler for verified webhook requests
app.post('/hubspot-webhook', async (req, res) => {
    console.log('Received and verified HubSpot webhook:', req.body);
    res.status(200).send('Webhook processed successfully');
});

app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});
