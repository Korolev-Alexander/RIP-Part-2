import React from 'react';
import { Link, useLocation } from 'react-router-dom';

interface BreadcrumbItem {
  label: string;
  path?: string;
}

const Breadcrumbs: React.FC = () => {
  const location = useLocation();
  
  const generateBreadcrumbs = (): BreadcrumbItem[] => {
    const paths = location.pathname.split('/').filter(path => path);
    
    const breadcrumbs: BreadcrumbItem[] = [
      { label: 'Главная', path: '/' }
    ];

    let currentPath = '';
    paths.forEach(path => {
      currentPath += `/${path}`;
      
      switch (path) {
        case 'devices':
          breadcrumbs.push({ label: 'Устройства', path: currentPath });
          break;
        case 'orders':
          breadcrumbs.push({ label: 'Заявки', path: currentPath });
          break;
        default:
          if (path.match(/^\d+$/)) {
            breadcrumbs.push({ label: 'Детали' });
          } else {
            breadcrumbs.push({ label: path.charAt(0).toUpperCase() + path.slice(1) });
          }
      }
    });

    return breadcrumbs;
  };

  const breadcrumbs = generateBreadcrumbs();

  return (
    <nav aria-label="breadcrumb">
      <ol className="breadcrumb">
        {breadcrumbs.map((breadcrumb, index) => (
          <li 
            key={index} 
            className={`breadcrumb-item ${index === breadcrumbs.length - 1 ? 'active' : ''}`}
          >
            {breadcrumb.path && index < breadcrumbs.length - 1 ? (
              <Link to={breadcrumb.path}>{breadcrumb.label}</Link>
            ) : (
              <span>{breadcrumb.label}</span>
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
};

export default Breadcrumbs;