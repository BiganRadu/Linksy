import React, { useState } from 'react';
import {useParams} from 'react-router-dom';
import { useEffect } from 'react';
import {useDeviceSelectors} from 'react-device-detect'

const ipToInt32 = (ip: string): number => {
	return ip.split('.').reduce((acc, octet) => (acc << 8) + parseInt(octet, 10), 0) >>> 0;
};

const fetchIpAndLocation = async () => {
	try {
		const response = await fetch('https://ipapi.co/json/');
		const data = await response.json();
		return data;
	} catch (error) {
		return null;
	}
};


const Redirecter: React.FC = () => {
    // Use useParams to access the dynamic segment from the route
    const { linkID } = useParams<{ linkID: string }>();
	const [redirectLink, setRedirectLink] = useState<string | null>(null);
	const [error, setError] = useState<any>(null);
	console.log(linkID)

	useEffect(() => {
		const fetchDataAndPost = async (device: string, os: string) => {
			let userData = await fetchIpAndLocation();
			if (!userData) {
				userData = {
					ip: '0.0.0.0',
					country_name: 'Unknown',
				}
			}

			const { ip, country_name } = userData;

			try {
				const response = await fetch(`http://localhost:3000/redirect?link_id=${linkID}`, {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						device: device,
						os: os,
						ip: ipToInt32(ip),
						country: country_name
					}),
				});

				const result = await response.json();
				if (response.ok) {
					setRedirectLink(result.redirect_link);
				} else {
					setError(result.error);
				}
			} catch (error) {
				console.error('Error during fetch:', error);
				setError(error);
			}
		};

		const [_, data] = useDeviceSelectors(window.navigator.userAgent);
		fetchDataAndPost(data.device.type != "" ? data.device.type : "unknown", data.os.name != "" ? data.os.name : "unknown");
	}, []);
	if (error) {
		return <div>Error: {error.message}</div>;
	}

	if (!redirectLink) {
		return <div>Loading...</div>;
	}
	        // Check if the URL starts with 'http://' or 'https://'
	if (!/^https?:\/\//i.test(redirectLink)) {
		setRedirectLink('https://' + redirectLink)
	}
			// Redirect to the formatted URL
	window.location.href = redirectLink;
}

export default Redirecter;