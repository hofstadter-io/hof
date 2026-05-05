export const prices = {
  "gemini-3-flash": {
    input: {
      cutoff: 0,
      short: 0.50,
      long: 0.50
    },
    cache: {
      cutoff: 0,
      short: .05,
      long: .05
    },
    output: {
      cutoff: 0,
      short: 3.00,
      long: 3.00
    }
  },
  "gemini-3-pro": {
    input: {
      cutoff: 200000,
      short: 2.00,
      long: 4.00
    },
    cache: {
      cutoff: 200000,
      short: .20,
      long: .40
    },
    output: {
      cutoff: 200000,
      short: 12.00,
      long: 18.00
    }
  }
}