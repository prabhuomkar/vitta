import React from 'react';
import { Box, Button, useTheme } from '@chakra-ui/react';
import { ArrowBackIcon, ArrowForwardIcon } from '@chakra-ui/icons';

const Pagination = ({
  page,
  totalPages,
  goToNextPage,
  goToPreviousPage,
  hasNextPage,
  loading
}) => {
  const theme = useTheme();
  const primaryColor = theme.colors.primary;

  return (
    <Box display="flex" justifyContent="right" alignItems="center" mt="4">
      <Button
        onClick={goToPreviousPage}
        colorScheme={primaryColor}
        bg={primaryColor}
        leftIcon={<ArrowBackIcon boxSize="14px" />}
        size="xs"
        hidden={page === 1}
        isDisabled={loading}
      >
        Previous
      </Button>
      <Box mx="2" fontSize="xs" display="inline-block">
        {totalPages > 0 && (
          <>
            Page {page} / {totalPages}
          </>
        )}
      </Box>
      <Button
        onClick={goToNextPage}
        colorScheme={primaryColor}
        bg={primaryColor}
        rightIcon={<ArrowForwardIcon boxSize="14px" />}
        size="xs"
        hidden={!hasNextPage}
        isDisabled={loading}
      >
        Next
      </Button>
    </Box>
  );
};

export default Pagination;
