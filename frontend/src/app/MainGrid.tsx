import * as React from 'react';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import AnalyticsRoundedIcon from '@mui/icons-material/AnalyticsRounded';
import LinkIcon from '@mui/icons-material/Link';
import QrCodeIcon from '@mui/icons-material/QrCode';
import Button from '@mui/material/Button';

export default function MainGrid() {

	const handleRedirect = (category: string) => {
		window.location.href = '/app/' + category;
	}
	
  return (
    <Box sx={{ width: '100%', maxWidth: '800px', paddingTop: 4}}>
      <Typography component="h2" variant="h3" sx={{ mb: 2 }}>
        Your connections platform
      </Typography>
      <Grid container spacing={2} marginTop={2}>
        <Grid item xs={4}>
          <Box sx={{ border: '1px solid', p: 2, display: 'flex', alignItems: 'center' }}>
            <LinkIcon sx={{ fontSize: 40, mr: 2 }} />
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
              <Typography variant="h6" component="div">
                Links
              </Typography>
              <Button sx={{ mt: 1, alignSelf: 'center' }} onClick={() => handleRedirect('links')}>View Links</Button>
            </Box>
          </Box>
        </Grid>
        <Grid item xs={4}>
          <Box sx={{ border: '1px solid', p: 2, display: 'flex', alignItems: 'center' }}>
            <QrCodeIcon sx={{ fontSize: 40, mr: 2 }} />
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
              <Typography variant="h6" component="div" align="center">
                QR Codes
              </Typography>
              <Button sx={{ mt: 1, alignSelf: 'center' }} onClick={() => handleRedirect('qrcodes')}>View QR Codes</Button>
            </Box>
          </Box>
        </Grid>
		<Grid item xs={4}>
          <Box sx={{ border: '1px solid', p: 2, display: 'flex', alignItems: 'center' }}>
            <AnalyticsRoundedIcon sx={{ fontSize: 40, mr: 2 }} />
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
              <Typography variant="h6" component="div">
                Analytics
              </Typography>
              <Button sx={{ mt: 1, alignSelf: 'center' }} onClick={() => handleRedirect('analytics')}>View Analytics</Button>
            </Box>
          </Box>
        </Grid>
      </Grid>
      <Typography component="h2" variant="h4" sx={{ mb: 2, mt: 3, paddingTop: 1 }}>
        Description
      </Typography>
	<Typography component="div" variant="h6" sx={{ mb: 2 }}>
		Our platform offers a comprehensive suite of tools to manage and analyze your connections:
	</Typography>
	<ul>
		<li>
			<Typography component="span" variant="h6">
				<strong>Links:</strong> Manage and organize your links efficiently.
			</Typography>
		</li>
		<li>
			<Typography component="span" variant="h6">
				<strong>QR Codes:</strong> Generate and track QR codes for easy sharing.
			</Typography>
		</li>
		<li>
			<Typography component="span" variant="h6">
				<strong>Analytics:</strong> Gain insights into your connection data with detailed analytics.
			</Typography>
		</li>
	</ul>
	<Typography component="div" variant="h6" sx={{ mt: 2 }}>
		Join us today and take control of your connections with our user-friendly platform.
	</Typography>
    </Box>
  );
}