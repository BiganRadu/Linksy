import React from 'react';
import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  IconButton,
  Link,
  Box
} from '@mui/material';
import ShareIcon from '@mui/icons-material/Share';
import CopyIcon from '@mui/icons-material/ContentCopy';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

interface LinkCardInfo {
	title: string;
	link_id: string;
	original_url: string;
}
const CustomCard: React.FC<LinkCardInfo> = ({ title, link_id, original_url }) => {
	const copyToClipboard = () => {
	  navigator.clipboard.writeText(`http://bit.ly/${link_id}`);
	};
  
	return (
	  <Card sx={{ width: '100%', maxWidth: 600, margin: 'auto', mt: 2}}>
		<CardContent>
		  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
			<Typography variant="h6" component="div">
			  {title}
			</Typography>
			<Box>
			  <Button variant="contained" size="small" startIcon={<CopyIcon />} onClick={copyToClipboard}>
				Copy
			  </Button>
			  <IconButton aria-label="share">
				<ShareIcon />
			  </IconButton>
			  <IconButton aria-label="edit">
				<EditIcon />
			  </IconButton>
			  <IconButton aria-label="delete">
				<DeleteIcon />
			  </IconButton>
			</Box>
		  </Box>
		  <Link href={`http://bit.ly/${link_id}`} target="_blank" sx={{ fontSize: 14, display: 'block', mb: 1, color:'blue'}}>
			bit.ly/{link_id}
		  </Link>
		  <Link href={original_url} target="_blank" sx={{ fontSize: 14, display: 'block' }}>
			{original_url}
		  </Link>
		  <Typography sx={{ fontSize: 12, mt: 2 }} color="text.secondary">
			Oct 13, 2024
		  </Typography>
		  <Typography sx={{ fontSize: 12 }} color="text.secondary">
			No tags
		  </Typography>
		</CardContent>
	  </Card>
	);
};
  

export default CustomCard;