describe('Página principal - Lista de Animales', () => {
  beforeEach(() => {
    cy.intercept('GET', 'http://localhost:3000/animals', { fixture: 'animals.json' }).as('getAnimals')
    cy.visit('/')
  })

  it('muestra la tabla con los animales', () => {
    cy.wait('@getAnimals')
    cy.get('table').should('exist')
    cy.get('table tbody tr').should('have.length.greaterThan', 0)
    cy.get('table').contains('Leo').should('exist')
  })
})

