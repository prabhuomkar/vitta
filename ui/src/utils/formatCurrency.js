const formatCurrency = amount => {
  if (amount === null || amount === undefined) return '₹0';

  const number = Number(amount);
  const formatted = number.toLocaleString('en-IN', {
    style: 'currency',
    currency: 'INR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  });

  return formatted;
};

export { formatCurrency };
