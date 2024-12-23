import { Box, Container, Typography } from '@mui/material';

function Home() {
  return (
    <Container maxWidth="lg">
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          py: 8
        }}
      >
        <Typography
          variant="h1"
          component="h1"
          sx={{
            fontSize: { xs: '2.5rem', md: '3.75rem' },
            fontWeight: 700,
            textAlign: 'center'
          }}
        >
          Home Page
        </Typography>
      </Box>
    </Container>
  );
}

export default Home;