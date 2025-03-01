import React from 'react';
import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  IconButton,
  Link,
  Box,
  Avatar
} from '@mui/material';
import ShareIcon from '@mui/icons-material/Share';
import CopyIcon from '@mui/icons-material/ContentCopy';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

interface LinkCardInfo {
	title: string;
	link_id: string;
	original_url: string;
	icon_url: string;
	created_at: number;
	onDelete?: () => void;
	onShare?: () => void;
	onEdit?: () => void;
}

const formatTimestamp = (timestamp: number): string => {
	const date = new Date(timestamp * 1000);
	return date.toLocaleDateString('en-US', {
	  month: 'short',
	  day: '2-digit',
	  year: 'numeric'
	});
};

const CustomCard: React.FC<LinkCardInfo> = ({ title, link_id, original_url, icon_url, created_at, onDelete, onShare, onEdit}) => {
	const copyToClipboard = () => {
	  navigator.clipboard.writeText(`http://bit.ly/${link_id}`);
	};
  
	return (
	  <Card sx={{ width: '100%', maxWidth: 600, margin: 'auto', mt: 2 }}>
		<CardContent>
		  <Box sx={{ display: 'flex', alignItems: 'center' }}>
			{/* Icon on the left */}
			<Link href={`http://bit.ly/${link_id}`} target="_blank" underline="none" sx={{ outline: 'none', "::before": { content: '""', display: 'none',} }}>
				<Avatar
					src={icon_url}
					alt="Icon"
					sx={{ width: 48, height: 48, mr: 2 }}
				/>
			</Link>

			{/* Card content */}
			<Box sx={{ flex: 1 }}>
			  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
				<Typography variant="h6" component="div">
				  {title}
				</Typography>
				<Box>
				  <Button variant="contained" size="small" startIcon={<CopyIcon />} onClick={copyToClipboard}>
					Copy
				  </Button>
				  <IconButton aria-label="share" onClick={onShare}>
					<ShareIcon />
				  </IconButton>
				  <IconButton aria-label="edit" onClick={onEdit}>
					<EditIcon />
				  </IconButton>
				  <IconButton aria-label="delete" onClick={onDelete}>
					<DeleteIcon />
				  </IconButton>
				</Box>
			  </Box>
			  <Link href={`http://bit.ly/${link_id}`} target="_blank" sx={{ fontSize: 14, display: 'block', mb: 1, color: 'blue' }}>
				bit.ly/{link_id}
			  </Link>
			  <Link href={original_url} target="_blank" sx={{ fontSize: 14, display: 'block' }}>
				{original_url}
			  </Link>
			  <Typography sx={{ fontSize: 12, mt: 2 }} color="text.secondary">
				{formatTimestamp(created_at)}
			  </Typography>
			</Box>
		  </Box>
		</CardContent>
	  </Card>
	);
};
  

export default CustomCard;